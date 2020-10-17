package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/urfave/cli/v2"
	pluralize "github.com/web-ridge/go-pluralize"
	gqlgen_sqlboiler "github.com/web-ridge/gqlgen-sqlboiler/v2"
)

var pluralizer *pluralize.Client //nolint:gochecknoglobals

func init() { //nolint:gochecknoinits
	pluralizer = pluralize.NewClient()
}

const (
	indent    = "\t"
	lineBreak = "\n"
)

// global configs

func main() {
	var modelDirectory string
	var outputFile string
	var mutations bool
	var batchUpdate bool
	var batchCreate bool
	var batchDelete bool
	var skipInputFields cli.StringSlice
	var directives cli.StringSlice
	var pagination string

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "input",
				Value:       "models",
				Usage:       "directory where the sqlboiler models are",
				Destination: &modelDirectory,
			},
			&cli.StringFlag{
				Name:        "output",
				Value:       "schema.graphql",
				Usage:       "filepath for schema",
				Destination: &outputFile,
			},
			&cli.StringSliceFlag{
				Name: "skip-input-fields",
				Usage: "input names which should be skipped: e.g. --skip-input-fields=userId --skip-input-fields=" +
					"organizationId --skip-input-fields=createdAt",
				Destination: &skipInputFields,
			},
			&cli.StringSliceFlag{
				Name:        "directives",
				Usage:       "directives which should be added after resolvers e.g. isAuthenticated",
				Destination: &directives,
			},
			&cli.BoolFlag{
				Name:        "mutations",
				Usage:       "generate mutations for models",
				Value:       true,
				Destination: &mutations,
			},
			&cli.BoolFlag{
				Name:        "batch-update",
				Usage:       "generate batch update for models",
				Value:       true,
				Destination: &batchUpdate,
			},
			&cli.BoolFlag{
				Name:        "batch-create",
				Usage:       "generate batch create for models",
				Value:       true,
				Destination: &batchCreate,
			},
			&cli.BoolFlag{
				Name:        "batch-delete",
				Usage:       "generate batch delete for models",
				Value:       true,
				Destination: &batchDelete,
			},
			&cli.StringFlag{
				Name:        "pagination",
				Usage:       "generate pagination support for models",
				Value:       "",
				Destination: &pagination,
			},
		},
		Action: func(c *cli.Context) error {
			// Generate schema based on config
			schema := getSchema(
				modelDirectory,
				mutations,
				batchUpdate,
				batchCreate,
				batchDelete,
				skipInputFields.Value(),
				directives.Value(),
				pagination,
			)

			// TODO: Write schema to the configured location
			if fileExists(outputFile) { //nolint:nestif
				baseFile := filenameWithoutExtension(outputFile) +
					"-empty" +
					getFilenameExtension(outputFile)

				newOutputFile := filenameWithoutExtension(outputFile) +
					"-new" +
					getFilenameExtension(outputFile)

				// remove previous files if exist
				_ = os.Remove(baseFile)
				_ = os.Remove(newOutputFile)

				if err := writeContentToFile(newOutputFile, schema); err != nil {
					return fmt.Errorf("could not write schema to disk: %v", err)
				}
				if err := formatFile(outputFile); err != nil {
					return fmt.Errorf("could not format with prettier %v: %v", outputFile, err)
				}
				if err := formatFile(newOutputFile); err != nil {
					return fmt.Errorf("could not format with prettier %v: %v", newOutputFile, err)
				}

				// Three way merging done based on this answer
				// https://stackoverflow.com/a/9123563/2508481

				// Empty file as base per the stackoverflow answer
				name := "touch"
				args := []string{baseFile}
				out, err := exec.Command(name, args...).Output()
				if err != nil {
					fmt.Println("Executing command failed: ", name, strings.Join(args, " "))
					return fmt.Errorf("merging failed %v: %v", err, out)
				}

				// Let's do the merge
				name = "git"
				args = []string{"merge-file", outputFile, baseFile, newOutputFile}
				out, err = exec.Command(name, args...).Output()
				if err != nil {
					fmt.Println("Executing command failed: ", name, strings.Join(args, " "))
					// remove base file
					_ = os.Remove(baseFile)
					return fmt.Errorf("merging failed or had conflicts %v: %v", err, out)
				}

				fmt.Println("Merging done without conflicts: ", out)

				// remove files

				_ = os.Remove(baseFile)
				_ = os.Remove(newOutputFile)

				// fmt.Printf("The date is %s\n", out)
			} else {
				fmt.Printf("Write schema of %v bytes to %v \n", len(schema), outputFile)
				if err := writeContentToFile(outputFile, schema); err != nil {
					fmt.Println("Could not write schema to disk: ", err)
				}
				return formatFile(outputFile)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getFilenameExtension(fn string) string {
	return path.Ext(fn)
}

func filenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func formatFile(filename string) error {
	name := "prettier"
	args := []string{filename, "--write"}

	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return fmt.Errorf("executing command: '%v %v' failed with: %v, output: %v", name, strings.Join(args, " "), err, out)
	}
	// fmt.Println(fmt.Sprintf("Formatting of %v done", filename))
	return nil
}

func writeContentToFile(filename string, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not write %v to disk: %v", filename, err)
	}

	// Close file if this functions returns early or at the end
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			fmt.Println("Error while closing file: ", closeErr)
		}
	}()

	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("could not write content to file %v: %v", filename, err)
	}

	return nil
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

const queryHelperStructs = `
input IDFilter {
	equalTo: ID
	notEqualTo: ID
	in: [ID!]
	notIn: [ID!]
}

input StringFilter {
	equalTo: String
	notEqualTo: String

	in: [String!]
	notIn: [String!]

	startWith: String
	notStartWith: String

	endWith: String
	notEndWith: String

	contain: String
	notContain: String

	startWithStrict: String # Camel sensitive
	notStartWithStrict: String # Camel sensitive

	endWithStrict: String # Camel sensitive
	notEndWithStrict: String # Camel sensitive

	containStrict: String # Camel sensitive
	notContainStrict: String # Camel sensitive
}

input IntFilter {
	equalTo: Int
	notEqualTo: Int
	lessThan: Int
	lessThanOrEqualTo: Int
	moreThan: Int
	moreThanOrEqualTo: Int
	in: [Int!]
	notIn: [Int!]
}

input FloatFilter {
	equalTo: Float
	notEqualTo: Float
	lessThan: Float
	lessThanOrEqualTo: Float
	moreThan: Float
	moreThanOrEqualTo: Float
	in: [Float!]
	notIn: [Float!]
}

input BooleanFilter {
	equalTo: Boolean
	notEqualTo: Boolean
}
`

type Model struct {
	Name   string
	Fields []*Field
	// Implements *string
}

type Field struct {
	Name             string
	RelationName     string // posts
	RelationType     string // Page, User, Post
	Type             string // String, ID, Integer
	FullType         string // e.g String! or if array [String!]
	RelationFullType string // [Posts!]
	FullTypeOptional string // e.g. String or if array [String]
	BoilerField      *gqlgen_sqlboiler.BoilerField
}

//nolint:gocognit,gocyclo // TODO: refactor this
func getSchema(
	modelDirectory string,
	mutations bool,
	batchUpdate bool,
	batchCreate bool,
	batchDelete bool,
	skipInputFields []string,
	directivesSlice []string,
	pagination string,
) string {
	var s strings.Builder

	// Parse models and their fields based on the sqlboiler model directory
	boilerModels := gqlgen_sqlboiler.GetBoilerModels(modelDirectory)
	models := boilerModelsToModels(boilerModels)

	var fullDirectives []string // nolint:prealloc
	for _, defaultDirective := range directivesSlice {
		fullDirectives = append(fullDirectives, "@"+defaultDirective)
		s.WriteString(fmt.Sprintf("directive @%v on FIELD_DEFINITION", defaultDirective))
		s.WriteString(lineBreak)
	}
	s.WriteString(lineBreak)

	joinedDirectives := strings.Join(fullDirectives, " ")
	// Create basic structs e.g.
	// type User {
	// 	firstName: String!
	// 	lastName: String
	// 	isProgrammer: Boolean!
	// 	organization: Organization!
	// }
	for _, model := range models {
		s.WriteString("type " + model.Name + " {")
		s.WriteString(lineBreak)
		for _, field := range model.Fields {
			// e.g we have foreign key from user to organization
			// organizationID is clutter in your scheme
			// you only want Organization and OrganizationID should be skipped
			if field.BoilerField.IsRelation {
				s.WriteString(indent + field.RelationName + ": " + field.RelationFullType)
				s.WriteString(lineBreak)
			} else {
				s.WriteString(indent + field.Name + ": " + field.FullType)
				s.WriteString(lineBreak)
			}
		}
		s.WriteString("}")
		s.WriteString(lineBreak)
		s.WriteString(lineBreak)
	}

	// Add helpers for filtering lists
	s.WriteString(queryHelperStructs)
	s.WriteString(lineBreak)

	// generate filter structs per model
	for _, model := range models {
		// Ignore some specified input fields

		// Generate a type safe grapql filter

		// Generate the base filter
		// type UserFilter {
		// 	search: String
		// 	where: UserWhere
		// }
		s.WriteString("input " + model.Name + "Filter {")
		s.WriteString(lineBreak)
		s.WriteString(indent + "search: String")
		s.WriteString(lineBreak)
		s.WriteString(indent + "where: " + model.Name + "Where")
		s.WriteString(lineBreak)
		s.WriteString("}")
		s.WriteString(lineBreak)
		s.WriteString(lineBreak)
		// Generate a pagination struct
		if pagination == "offset" {
			// type UserPagination {
			// 	limit: Int!
			// 	page: Int!
			// }
			s.WriteString("input " + model.Name + "Pagination {")
			s.WriteString(lineBreak)
			s.WriteString(indent + "limit: Int!")
			s.WriteString(lineBreak)
			s.WriteString(indent + "page: Int!")
			s.WriteString(lineBreak)
			s.WriteString("}")
			s.WriteString(lineBreak)
			s.WriteString(lineBreak)
		}
		// Generate a where struct
		// type UserWhere {
		// 	id: IDFilter
		// 	title: StringFilter
		// 	organization: OrganizationWhere
		// 	or: FlowBlockWhere
		// 	and: FlowBlockWhere
		// }
		s.WriteString("input " + model.Name + "Where {")
		s.WriteString(lineBreak)
		for _, field := range model.Fields {
			if field.BoilerField.IsRelation {
				// Support filtering in relationships (atleast schema wise)
				s.WriteString(indent + field.RelationName + ": " + field.RelationType + "Where")
				s.WriteString(lineBreak)
			} else {
				s.WriteString(indent + field.Name + ": " + field.Type + "Filter")
				s.WriteString(lineBreak)
			}
		}
		s.WriteString(indent + "or: " + model.Name + "Where")
		s.WriteString(lineBreak)

		s.WriteString(indent + "and: " + model.Name + "Where")
		s.WriteString(lineBreak)

		s.WriteString("}")
		s.WriteString(lineBreak)
		s.WriteString(lineBreak)
	}

	s.WriteString("type Query {")
	s.WriteString(lineBreak)
	for _, model := range models {
		// single models
		s.WriteString(indent)
		s.WriteString(strcase.ToLowerCamel(model.Name) + "(id: ID!)")
		s.WriteString(": ")
		s.WriteString(model.Name + "!")
		s.WriteString(joinedDirectives)
		s.WriteString(lineBreak)

		// lists
		modelPluralName := pluralizer.Plural(model.Name)
		s.WriteString(indent)
		var paginationParameter string
		if pagination == "offset" {
			paginationParameter = ", pagination: " + model.Name + "Pagination"
		}
		s.WriteString(strcase.ToLowerCamel(modelPluralName) + "(filter: " + model.Name + "Filter" +
			paginationParameter + ")")
		s.WriteString(": ")
		s.WriteString("[" + model.Name + "!]!")
		s.WriteString(joinedDirectives)
		s.WriteString(lineBreak)
	}
	s.WriteString("}")
	s.WriteString(lineBreak)
	s.WriteString(lineBreak)

	// Generate input and payloads for mutatations
	if mutations { //nolint:nestif
		for _, model := range models {
			filteredFields := fieldsWithout(model.Fields, skipInputFields)

			modelPluralName := pluralizer.Plural(model.Name)
			// input UserCreateInput {
			// 	firstName: String!
			// 	lastName: String
			//	organizationId: ID!
			// }
			s.WriteString("input " + model.Name + "CreateInput {")
			s.WriteString(lineBreak)
			for _, field := range filteredFields {
				// id is not required in create and will be specified in update resolver
				if field.Name == "id" {
					continue
				}

				// not possible yet in input
				// TODO: make this possible for one-to-one structs?
				// only for foreign keys inside model itself
				if field.BoilerField.IsRelation && field.BoilerField.IsArray ||
					field.BoilerField.IsRelation && !strings.HasSuffix(field.BoilerField.Name, "ID") {
					continue
				}

				s.WriteString(indent + field.Name + ": " + field.FullType)
				s.WriteString(lineBreak)
			}
			s.WriteString("}")
			s.WriteString(lineBreak)
			s.WriteString(lineBreak)

			// input UserUpdateInput {
			// 	firstName: String!
			// 	lastName: String
			//	organizationId: ID!
			// }
			s.WriteString("input " + model.Name + "UpdateInput {")
			s.WriteString(lineBreak)
			for _, field := range filteredFields {
				// id is not required in create and will be specified in update resolver
				if field.Name == "id" {
					continue
				}
				// not possible yet in input
				// TODO: make this possible for one-to-one structs?
				// only for foreign keys inside model itself
				if field.BoilerField.IsRelation && field.BoilerField.IsArray ||
					field.BoilerField.IsRelation && !strings.HasSuffix(field.BoilerField.Name, "ID") {
					continue
				}

				s.WriteString(indent + field.Name + ": " + field.FullTypeOptional)
				s.WriteString(lineBreak)
			}
			s.WriteString("}")
			s.WriteString(lineBreak)
			s.WriteString(lineBreak)

			if batchCreate {
				s.WriteString("input " + modelPluralName + "CreateInput {")
				s.WriteString(lineBreak)
				s.WriteString(indent + strcase.ToLowerCamel(modelPluralName) + ": [" + model.Name + "CreateInput!]!")
				s.WriteString("}")
				s.WriteString(lineBreak)
				s.WriteString(lineBreak)
			}

			// if batchUpdate {
			// 	s.WriteString("input " + modelPluralName + "UpdateInput {")
			// 	s.WriteString(lineBreak)
			// 	s.WriteString(indent + strcase.ToLowerCamel(modelPluralName) + ": [" + model.Name + "UpdateInput!]!")
			// 	s.WriteString("}")
			// 	s.WriteString(lineBreak)
			// 	s.WriteString(lineBreak)
			// }

			// type UserPayload {
			// 	user: User!
			// }
			s.WriteString("type " + model.Name + "Payload {")
			s.WriteString(lineBreak)
			s.WriteString(indent + strcase.ToLowerCamel(model.Name) + ": " + model.Name + "!")
			s.WriteString(lineBreak)
			s.WriteString("}")
			s.WriteString(lineBreak)
			s.WriteString(lineBreak)

			// TODO batch, delete input and payloads

			// type UserDeletePayload {
			// 	id: ID!
			// }
			s.WriteString("type " + model.Name + "DeletePayload {")
			s.WriteString(lineBreak)
			s.WriteString(indent + "id: ID!")
			s.WriteString(lineBreak)
			s.WriteString("}")
			s.WriteString(lineBreak)
			s.WriteString(lineBreak)

			// type UsersPayload {
			// 	ids: [ID!]!
			// }
			if batchCreate {
				s.WriteString("type " + modelPluralName + "Payload {")
				s.WriteString(lineBreak)
				s.WriteString(indent + strcase.ToLowerCamel(modelPluralName) + ": [" + model.Name + "!]!")
				s.WriteString(lineBreak)
				s.WriteString("}")
				s.WriteString(lineBreak)
				s.WriteString(lineBreak)
			}

			// type UsersDeletePayload {
			// 	ids: [ID!]!
			// }
			if batchDelete {
				s.WriteString("type " + modelPluralName + "DeletePayload {")
				s.WriteString(lineBreak)
				s.WriteString(indent + "ids: [ID!]!")
				s.WriteString(lineBreak)
				s.WriteString("}")
				s.WriteString(lineBreak)
				s.WriteString(lineBreak)
			}
			// type UsersUpdatePayload {
			// 	ok: Boolean!
			// }
			if batchUpdate {
				s.WriteString("type " + modelPluralName + "UpdatePayload {")
				s.WriteString(lineBreak)
				s.WriteString(indent + "ok: Boolean!")
				s.WriteString(lineBreak)
				s.WriteString("}")
				s.WriteString(lineBreak)
				s.WriteString(lineBreak)
			}
		}

		// Generate mutation queries
		s.WriteString("type Mutation {")
		s.WriteString(lineBreak)
		for _, model := range models {
			modelPluralName := pluralizer.Plural(model.Name)

			// create single
			// e.g createUser(input: UserInput!): UserPayload!
			s.WriteString(indent)
			s.WriteString("create" + model.Name + "(input: " + model.Name + "CreateInput!)")
			s.WriteString(": ")
			s.WriteString(model.Name + "Payload!")
			s.WriteString(joinedDirectives)
			s.WriteString(lineBreak)

			// create multiple
			// e.g createUsers(input: [UsersInput!]!): UsersPayload!
			if batchCreate {
				s.WriteString(indent)
				s.WriteString("create" + modelPluralName + "(input: " + modelPluralName + "CreateInput!)")
				s.WriteString(": ")
				s.WriteString(modelPluralName + "Payload!")
				s.WriteString(joinedDirectives)
				s.WriteString(lineBreak)
			}

			// update single
			// e.g updateUser(id: ID!, input: UserInput!): UserPayload!
			s.WriteString(indent)
			s.WriteString("update" + model.Name + "(id: ID!, input: " + model.Name + "UpdateInput!)")
			s.WriteString(": ")
			s.WriteString(model.Name + "Payload!")
			s.WriteString(joinedDirectives)
			s.WriteString(lineBreak)

			// update multiple (batch update)
			// e.g updateUsers(filter: UserFilter, input: UsersInput!): UsersPayload!
			if batchUpdate {
				s.WriteString(indent)
				s.WriteString("update" + modelPluralName + "(filter: " + model.Name + "Filter, input: " +
					model.Name + "UpdateInput!)")
				s.WriteString(": ")
				s.WriteString(modelPluralName + "UpdatePayload!")
				s.WriteString(joinedDirectives)
				s.WriteString(lineBreak)
			}

			// delete single
			// e.g deleteUser(id: ID!): UserPayload!
			s.WriteString(indent)
			s.WriteString("delete" + model.Name + "(id: ID!)")
			s.WriteString(": ")
			s.WriteString(model.Name + "DeletePayload!")
			s.WriteString(joinedDirectives)
			s.WriteString(lineBreak)

			// delete multiple
			// e.g deleteUsers(filter: UserFilter, input: [UsersInput!]!): UsersPayload!
			if batchDelete {
				s.WriteString(indent)
				s.WriteString("delete" + modelPluralName + "(filter: " + model.Name + "Filter)")
				s.WriteString(": ")
				s.WriteString(modelPluralName + "DeletePayload!")
				s.WriteString(joinedDirectives)
				s.WriteString(lineBreak)
			}
		}
		s.WriteString("}")
		s.WriteString(lineBreak)
		s.WriteString(lineBreak)
	}

	return s.String()
}

func getFullType(fieldType string, isArray bool, isRequired bool) string {
	gType := fieldType

	if isArray {
		// To use a list type, surround the type in square brackets, so [Int] is a list of integers.
		gType = "[" + gType + "]"
	}
	if isRequired {
		// Use an exclamation point to indicate a type cannot be nullable,
		// so String! is a non-nullable string.
		gType += "!"
	}
	return gType
}

func boilerModelsToModels(boilerModels []*gqlgen_sqlboiler.BoilerModel) []*Model {
	models := make([]*Model, len(boilerModels))
	for i, boilerModel := range boilerModels {
		models[i] = &Model{
			Name:   boilerModel.Name,
			Fields: boilerFieldsToFields(boilerModel.Fields),
		}
	}
	return models
}

func boilerFieldsToFields(boilerFields []*gqlgen_sqlboiler.BoilerField) []*Field {
	fields := make([]*Field, len(boilerFields))
	for i, boilerField := range boilerFields {
		fields[i] = boilerFieldToField(boilerField)
	}
	return fields
}

func boilerFieldToField(boilerField *gqlgen_sqlboiler.BoilerField) *Field {
	var relationName string
	var relationType string
	var relationFullType string
	if boilerField.Relationship != nil {
		relationName = strcase.ToLowerCamel(boilerField.RelationshipName)
		relationType = boilerField.Relationship.Name

		relationFullType = getFullType(
			relationType,
			boilerField.IsArray,
			boilerField.IsRequired,
		)
	}

	t := toGraphQLType(boilerField.Name, boilerField.Type)
	return &Field{
		Name:             toGraphQLName(boilerField.Name),
		RelationName:     relationName,
		RelationType:     relationType,
		Type:             t,
		FullType:         getFullType(t, boilerField.IsArray, boilerField.IsRequired),
		FullTypeOptional: getFullType(t, boilerField.IsArray, false),
		RelationFullType: relationFullType,
		BoilerField:      boilerField,
	}
}

func toGraphQLName(fieldName string) string {
	graphqlName := fieldName

	// Golang ID to Id the right way
	// Primary key
	if graphqlName == "ID" {
		graphqlName = "id"
	}

	if graphqlName == "URL" {
		graphqlName = "url"
	}

	// e.g. OrganizationID, TODO: more robust solution?
	graphqlName = strings.Replace(graphqlName, "ID", "Id", -1)
	graphqlName = strings.Replace(graphqlName, "URL", "Url", -1)

	return strcase.ToLowerCamel(graphqlName)
}

func toGraphQLType(fieldName, boilerType string) string {
	lowerFieldName := strings.ToLower(fieldName)
	lowerBoilerType := strings.ToLower(boilerType)

	if strings.HasSuffix(lowerFieldName, "id") {
		return "ID"
	}
	if strings.Contains(lowerBoilerType, "string") {
		return "String"
	}
	if strings.Contains(lowerBoilerType, "int") {
		return "Int"
	}
	if strings.Contains(lowerBoilerType, "decimal") || strings.Contains(lowerBoilerType, "float") {
		return "Float"
	}
	if strings.Contains(lowerBoilerType, "bool") {
		return "Boolean"
	}

	// TODO: make this a scalar or something configurable?
	// I like to use unix here
	if strings.Contains(lowerBoilerType, "time") {
		return "Int"
	}

	// E.g. UserSlice
	boilerType = strings.TrimSuffix(boilerType, "Slice")

	return boilerType
}

func fieldsWithout(fields []*Field, skipFieldNames []string) []*Field {
	var filteredFields []*Field
	for _, field := range fields {
		if !sliceContains(skipFieldNames, field.Name) {
			filteredFields = append(filteredFields, field)
		}
	}
	return filteredFields
}

func sliceContains(slice []string, v string) bool {
	for _, s := range slice {
		if s == v {
			return true
		}
	}
	return false
}
