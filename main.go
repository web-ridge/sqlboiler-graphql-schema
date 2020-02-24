package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"

	pluralize "github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/micro/cli"
	"github.com/web-ridge/gqlgen-sqlboiler/boiler"
)

var pluralizer *pluralize.Client

func init() {
	pluralizer = pluralize.NewClient()
}

const indent = "\t"

// global configs

func main() {
	var modelDirectory string
	var outputFile string
	var mutations bool
	var batchUpdate bool
	var batchCreate bool
	var batchDelete bool
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
			&cli.BoolTFlag{
				Name:        "mutations",
				Usage:       "generate mutations for models",
				Destination: &mutations,
			},
			&cli.BoolTFlag{
				Name:        "batch-update",
				Usage:       "generate batch update for models",
				Destination: &batchUpdate,
			},
			&cli.BoolTFlag{
				Name:        "batch-create",
				Usage:       "generate batch create for models",
				Destination: &batchCreate,
			},
			&cli.BoolTFlag{
				Name:        "batch-delete",
				Usage:       "generate batch delete for models",
				Destination: &batchDelete,
			},
		},
		Action: func(c *cli.Context) error {

			// Generate schema based on config
			schema := getSchema(modelDirectory, mutations, batchUpdate, batchCreate, batchDelete)

			// TODO: Write schema to the configured location
			if fileExists(outputFile) {

				baseFile := filenameWithoutExtension(outputFile) +
					"-empty" +
					getFilenameExtension(outputFile)

				newOutputFile := filenameWithoutExtension(outputFile) +
					"-new" +
					getFilenameExtension(outputFile)

				// remove previous files if exist
				os.Remove(baseFile)
				os.Remove(newOutputFile)

				if err := writeContentToFile(newOutputFile, schema); err != nil {
					return fmt.Errorf("Could not write schema to disk: %v", err)
				}
				if err := formatFile(outputFile); err != nil {
					return fmt.Errorf("Could not format with prettier %v: %v", outputFile, err)
				}
				if err := formatFile(newOutputFile); err != nil {
					return fmt.Errorf("Could not format with prettier %v: %v", newOutputFile, err)
				}

				// Three way merging done based on this answer
				// https://stackoverflow.com/a/9123563/2508481

				// Empty file as base per the stackoverflow answer
				name := "touch"
				args := []string{baseFile}
				out, err := exec.Command(name, args...).Output()
				if err != nil {
					fmt.Println("Executing command failed: ", name, strings.Join(args, " "))
					return fmt.Errorf("Merging failed %v: %v", err, out)
				}

				// Let's do the merge
				name = "git"
				args = []string{"merge-file", outputFile, baseFile, newOutputFile}
				out, err = exec.Command(name, args...).Output()
				if err != nil {
					fmt.Println("Executing command failed: ", name, strings.Join(args, " "))
					// remove base file
					os.Remove(baseFile)
					return fmt.Errorf("Merging failed or had conflicts %v: %v", err, out)
				}

				fmt.Println("Merging done without conflicts: ", out)

				// remove files
				os.Remove(baseFile)
				os.Remove(newOutputFile)

				// fmt.Printf("The date is %s\n", out)

			} else {
				fmt.Println(fmt.Sprintf("Write schema of %v bytes to %v", len(schema), outputFile))
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
		return fmt.Errorf("Executing command: '%v %v' failed with: %v, output: %v", name, strings.Join(args, " "), err, out)
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
	in: [ID!]
	notIn: [ID!]
}

input StringFilter {
	equalTo: String
	in: [String!]
	notIn: [String!]

	startsWith: String
	endsWith: String
	contains: String

	startsWithStrict: String # Camel sensitive
	endsWithStrict: String # Camel sensitive
	containsStrict: String # Camel sensitive
}

input IntFilter {
	equalTo: Int
	lessThan: Int
	lessThanOrEqualTo: Int
	moreThan: Int
	moreThanOrEqualTo: Int
	in: [Int!]
	notIn: [Int!]
}

input FloatFilter {
	equalTo: Float
	lessThan: Float
	lessThanOrEqualTo: Float
	moreThan: Float
	moreThanOrEqualTo: Float
	in: [Float!]
	notIn: [Float!]
}

input BooleanFilter {
	equalTo: Boolean
}
`

type Model struct {
	Name   string
	Fields []*Field
	// Implements *string
}
type Field struct {
	Name       string
	Type       string // String, ID, Integer
	FullType   string // e.g String! or if array [String!]
	BoilerName string
	BoilerType string
	IsRequired bool
	IsArray    bool
	IsRelation bool
}

type BoilerType struct {
	Name string
	Type string
}

func getSchema(
	modelDirectory string,
	mutations bool,
	batchUpdate bool,
	batchCreate bool,
	batchDelete bool,
) string {
	var schema strings.Builder

	// Parse models and their fields based on the sqlboiler model directory
	models := parseModelsAndFieldsFromBoiler(getSortedBoilerTypes(modelDirectory))

	// Create basic structs e.g.
	// type User {
	// 	firstName: String!
	// 	lastName: String
	// 	isProgrammer: Boolean!
	// 	organization: Organization!
	// }
	for _, model := range models {

		schema.WriteString("type " + model.Name + " {")
		schema.WriteString("\n")
		for _, field := range model.Fields {
			// e.g we have foreign key from user to organization
			// organizationID is clutter in your scheme
			// you only want Organization and OrganizationID should be skipped
			if !relationExist(field, model.Fields) {
				schema.WriteString(indent + field.Name + ": " + field.FullType)
				schema.WriteString("\n")
			}

		}
		schema.WriteString("}")
		schema.WriteString("\n")
		schema.WriteString("\n")
	}

	// Add helpers for filtering lists
	schema.WriteString(queryHelperStructs)
	schema.WriteString("\n")

	// generate filter structs per model
	for _, model := range models {

		// Generate a type safe grapql filter

		// Generate the base filter
		// type UserFilter {
		// 	search: String
		// 	where: UserWhere
		// }
		schema.WriteString("input " + model.Name + "Filter {")
		schema.WriteString("\n")
		schema.WriteString(indent + "search: String")
		schema.WriteString("\n")
		schema.WriteString(indent + "where: " + model.Name + "Where")
		schema.WriteString("\n")
		schema.WriteString("}")
		schema.WriteString("\n")
		schema.WriteString("\n")
		// Generate a where struct
		// type UserWhere {
		// 	id: IDFilter
		// 	title: StringFilter
		// 	organization: OrganizationWhere
		// 	or: FlowBlockWhere
		// 	and: FlowBlockWhere
		// }
		schema.WriteString("input " + model.Name + "Where {")
		schema.WriteString("\n")
		for _, field := range model.Fields {
			if field.IsRelation {
				// Support filtering in relationships (atleast schema wise)
				schema.WriteString(indent + field.Name + ": " + field.Type + "Where")
				schema.WriteString("\n")
			} else {
				schema.WriteString(indent + field.Name + ": " + field.Type + "Filter")
				schema.WriteString("\n")
			}
		}
		schema.WriteString(indent + "or: " + model.Name + "Where")
		schema.WriteString("\n")

		schema.WriteString(indent + "and: " + model.Name + "Where")
		schema.WriteString("\n")

		schema.WriteString("}")
		schema.WriteString("\n")
		schema.WriteString("\n")
	}

	schema.WriteString("type Query {")
	schema.WriteString("\n")
	for _, model := range models {
		// single models
		schema.WriteString(indent)
		schema.WriteString(strcase.ToLowerCamel(model.Name) + "(id: ID!)")
		schema.WriteString(": ")
		schema.WriteString(model.Name + "!")
		schema.WriteString("\n")

		// lists
		modelArray := pluralizer.Plural(model.Name)
		schema.WriteString(indent)
		// TODO: pagination
		schema.WriteString(strcase.ToLowerCamel(modelArray) + "(filter: " + model.Name + "Filter)")
		schema.WriteString(": ")
		schema.WriteString("[" + model.Name + "!]!")
		schema.WriteString("\n")

	}
	schema.WriteString("}")
	schema.WriteString("\n")
	schema.WriteString("\n")

	// Generate input and payloads for mutatations
	if mutations {
		for _, model := range models {

			modelArray := pluralizer.Plural(model.Name)
			// input UserInput {
			// 	firstName: String!
			// 	lastName: String
			//	organizationId: ID!
			// }
			schema.WriteString("input " + model.Name + "Input {")
			schema.WriteString("\n")
			for _, field := range model.Fields {
				// id is not required in create and will be specified in update resolver
				if field.Name == "id" {
					continue
				}
				// not possible yet in input
				if field.IsRelation {
					continue
				}
				schema.WriteString(indent + field.Name + ": " + field.FullType)
				schema.WriteString("\n")
			}
			schema.WriteString("}")
			schema.WriteString("\n")
			schema.WriteString("\n")

			if batchCreate {
				schema.WriteString("input " + modelArray + "Input {")
				schema.WriteString("\n")
				schema.WriteString(indent + strcase.ToLowerCamel(modelArray) + ": [" + model.Name + "Input!]!")
				schema.WriteString("}")
				schema.WriteString("\n")
				schema.WriteString("\n")
			}

			// type UserPayload {
			// 	user: User!
			// }
			schema.WriteString("type " + model.Name + "Payload {")
			schema.WriteString("\n")
			schema.WriteString(indent + strcase.ToLowerCamel(model.Name) + ": " + model.Name + "!")
			schema.WriteString("\n")
			schema.WriteString("}")
			schema.WriteString("\n")
			schema.WriteString("\n")

			// TODO batch, delete input and payloads

			// type UserDeletePayload {
			// 	id: ID!
			// }
			schema.WriteString("type " + model.Name + "DeletePayload {")
			schema.WriteString("\n")
			schema.WriteString(indent + "id: ID!")
			schema.WriteString("\n")
			schema.WriteString("}")
			schema.WriteString("\n")
			schema.WriteString("\n")

			// type UsersPayload {
			// 	ids: [ID!]!
			// }
			if batchCreate {
				schema.WriteString("type " + modelArray + "Payload {")
				schema.WriteString("\n")
				schema.WriteString(indent + strcase.ToLowerCamel(modelArray) + ": [" + model.Name + "!]!")
				schema.WriteString("\n")
				schema.WriteString("}")
				schema.WriteString("\n")
				schema.WriteString("\n")
			}

			// type UsersDeletePayload {
			// 	ids: [ID!]!
			// }
			if batchDelete {
				schema.WriteString("type " + modelArray + "DeletePayload {")
				schema.WriteString("\n")
				schema.WriteString(indent + "ids: [ID!]!")
				schema.WriteString("\n")
				schema.WriteString("}")
				schema.WriteString("\n")
				schema.WriteString("\n")
			}
			// type UsersUpdatePayload {
			// 	ok: Boolean!
			// }
			if batchUpdate {
				schema.WriteString("type " + modelArray + "UpdatePayload {")
				schema.WriteString("\n")
				schema.WriteString(indent + "ok: Boolean!")
				schema.WriteString("\n")
				schema.WriteString("}")
				schema.WriteString("\n")
				schema.WriteString("\n")
			}

		}
	}

	// Generate mutation queries

	if mutations {
		schema.WriteString("type Mutation {")
		schema.WriteString("\n")
		for _, model := range models {

			modelArray := pluralizer.Plural(model.Name)

			// create single
			// e.g createUser(input: UserInput!): UserPayload!
			schema.WriteString(indent)
			schema.WriteString("create" + model.Name + "(input: " + model.Name + "Input!)")
			schema.WriteString(": ")
			schema.WriteString(model.Name + "Payload!")
			schema.WriteString("\n")

			// create multiple
			// e.g createUsers(input: [UsersInput!]!): UsersPayload!
			if batchCreate {
				schema.WriteString(indent)
				schema.WriteString("create" + modelArray + "(input: " + modelArray + "Input!)")
				schema.WriteString(": ")
				schema.WriteString(modelArray + "Payload!")
				schema.WriteString("\n")
			}

			// update single
			// e.g updateUser(id: ID!, input: UserInput!): UserPayload!
			schema.WriteString(indent)
			schema.WriteString("update" + model.Name + "(id: ID!, input: " + model.Name + "Input!)")
			schema.WriteString(": ")
			schema.WriteString(model.Name + "Payload!")
			schema.WriteString("\n")

			// update multiple (batch update)
			// e.g updateUsers(filter: UserFilter, input: [UsersInput!]!): UsersPayload!
			if batchUpdate {
				schema.WriteString(indent)
				schema.WriteString("update" + modelArray + "(filter: " + model.Name + "Filter, input: " + modelArray + "Input!)")
				schema.WriteString(": ")
				schema.WriteString(modelArray + "UpdatePayload!")
				schema.WriteString("\n")
			}

			// delete single
			// e.g deleteUser(id: ID!): UserPayload!
			schema.WriteString(indent)
			schema.WriteString("delete" + model.Name + "(id: ID!)")
			schema.WriteString(": ")
			schema.WriteString(model.Name + "DeletePayload!")
			schema.WriteString("\n")

			// delete multiple
			// e.g deleteUsers(filter: UserFilter, input: [UsersInput!]!): UsersPayload!
			if batchDelete {
				schema.WriteString(indent)
				schema.WriteString("delete" + modelArray + "(filter: " + model.Name + "Filter)")
				schema.WriteString(": ")
				schema.WriteString(modelArray + "DeletePayload!")
				schema.WriteString("\n")
			}

		}
		schema.WriteString("}")
		schema.WriteString("\n")
		schema.WriteString("\n")
	}

	return schema.String()
}

// parseModelsAndFieldsFromBoiler since these are like User.ID, User.Organization and we want them grouped by
// modelName and their belonging fields.
func parseModelsAndFieldsFromBoiler(boilerTypes []*BoilerType) []*Model {

	// sortedModelNames is needed to get the right order back of the models since we want the same order every time
	// this program has ran.
	modelNames := []string{}

	// fieldsPerModelName is needed to group the fields per model, so we can get all fields per modelName later on
	fieldsPerModelName := map[string][]*Field{}

	// Anonymous function because this is used 2 times it prevents duplicated code
	// It's automatically inits an empty field array if it does not exist yet
	var addFieldsToModel = func(modelName string, field *Field) {
		modelNames = appendIfMissing(modelNames, modelName)
		_, ok := fieldsPerModelName[modelName]
		if !ok {
			fieldsPerModelName[modelName] = []*Field{}
		}
		fieldsPerModelName[modelName] = append(fieldsPerModelName[modelName], field)
	}

	// Let's parse boilerTypes to models and fields
	for _, boiler := range boilerTypes {

		// split on . input is like e.g. User.ID
		splitted := strings.Split(boiler.Name, ".")
		// result in e.g. User
		modelName := splitted[0]
		// result in e.g. ID
		boilerFieldName := splitted[1]

		// handle names with lowercase e.g. userR, userL or other sqlboiler extra's
		if isFirstCharacterLowerCase(modelName) {

			// It's the relations of the model
			// let's add them so we can use them later
			if strings.HasSuffix(modelName, "R") {
				modelName = strcase.ToCamel(strings.TrimSuffix(modelName, "R"))
				relationField := toField(boilerFieldName, boiler.Type)
				relationField.IsRelation = true
				addFieldsToModel(modelName, relationField)
			}

			// ignore the default handling since this field is already handled
			continue
		}

		// Ignore these since these are sqlboiler helper structs for preloading relationships
		if boilerFieldName == "L" || boilerFieldName == "R" {
			continue
		}

		addFieldsToModel(modelName, toField(boilerFieldName, boiler.Type))
	}
	sort.Strings(modelNames)

	// Let's generate the models in the same order as the sqlboiler structs were parsed
	models := make([]*Model, len(modelNames))
	for i, modelName := range modelNames {
		fields := fieldsPerModelName[modelName]

		// check if required based on foreign keys
		for _, f := range fields {
			if f.IsRelation {
				f.IsRequired = foreignKeyIsRequired(f, fields)
				f.FullType = getFullType(f.Type, f.IsArray, f.IsRequired)
			}
		}
		// for _, f := range fields {
		// 	if f.IsRelation {
		// 		fmt.Println("Is", modelName, f.BoilerName, "required?")
		// 		fmt.Println(f.IsRequired)
		// 	}
		// }

		models[i] = &Model{
			Name:   modelName,
			Fields: fields,
		}
	}
	return models
}

// getSortedBoilerTypes orders the sqlboiler struct in an ordered slice of BoilerType
func getSortedBoilerTypes(modelDirectory string) (sortedBoilerTypes []*BoilerType) {
	boilerTypeMap, _, boilerTypeOrder := boiler.ParseBoilerFile(modelDirectory)

	boilerTypeKeys := make([]string, 0, len(boilerTypeMap))
	for k := range boilerTypeMap {
		boilerTypeKeys = append(boilerTypeKeys, k)
	}

	// order same way as sqlboiler fields with one exception
	// let createdAt, updatedAt and deletedAt as last
	sort.Slice(boilerTypeKeys, func(i, b int) bool {

		aKey := boilerTypeKeys[i]
		bKey := boilerTypeKeys[b]

		aOrder := boilerTypeOrder[aKey]
		bOrder := boilerTypeOrder[bKey]

		higherOrders := []string{"createdat", "updatedat", "deletedat"}
		for i, higherOrder := range higherOrders {
			if strings.HasSuffix(strings.ToLower(aKey), higherOrder) {
				aOrder += 1000000 + i
			}
			if strings.HasSuffix(strings.ToLower(bKey), higherOrder) {
				bOrder += 10000000 + i
			}
		}

		return aOrder < bOrder
	})

	for _, modelAndField := range boilerTypeKeys {
		// fmt.Println(modelAndField)
		sortedBoilerTypes = append(sortedBoilerTypes, &BoilerType{
			Name: modelAndField,
			Type: boilerTypeMap[modelAndField],
		})
	}
	return
}

func foreignKeyIsRequired(relation *Field, foreignKeys []*Field) bool {
	findKey := relation.BoilerName + "ID"
	for _, foreignKey := range foreignKeys {
		if foreignKey.BoilerName == findKey {
			return isRequired(foreignKey.BoilerType)
		}
	}
	return false
}
func relationExist(field *Field, fields []*Field) bool {
	// e.g we have foreign key from user to organization
	// organizationID is clutter in your scheme
	// you only want Organization and OrganizationID should be skipped

	// ID can't possible be a relationship
	if field.BoilerName == "ID" {
		return false
	}

	// Ok, if it ends on ID it could have a relation ship
	if strings.HasSuffix(field.BoilerName, "ID") {
		for _, checkWithField := range fields {
			// don't compare to itself
			if field == checkWithField {
				continue
			}

			if checkWithField.BoilerName == strings.TrimSuffix(field.BoilerName, "ID") && checkWithField.IsRelation {
				// fmt.Println("Remove from fields")
				// fmt.Println(checkWithField.BoilerName)
				// fmt.Println(strings.TrimSuffix(field.BoilerName, "ID"))
				return true
			}
		}

	}

	return false
}

func isRequired(boilerType string) bool {
	if strings.HasPrefix(boilerType, "null.") || strings.HasPrefix(boilerType, "*") {
		return false
	}
	return true
}

func isArray(boilerType string) bool {
	return strings.HasSuffix(boilerType, "Slice")
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
		gType = gType + "!"
	}
	return gType
}

func toField(boilerName, boilerType string) *Field {
	t := toGraphQLType(boilerName, boilerType)
	array := isArray(boilerType)
	required := isRequired(boilerType)
	return &Field{
		Name:       toGraphQLName(boilerName, boilerType),
		Type:       t,
		FullType:   getFullType(t, array, required),
		BoilerName: boilerName,
		BoilerType: boilerType,
		IsRequired: required,
		IsArray:    array,
	}
}

func toGraphQLName(fieldName, boilerType string) string {
	graphqlName := fieldName

	// Golang ID to Id the right way
	// Primary key
	if graphqlName == "ID" {
		graphqlName = "id"
	}

	if graphqlName == "URL" {
		graphqlName = "url"
	}

	// e.g. OrganizationID
	graphqlName = strings.Replace(graphqlName, "ID", "Id", -1)
	graphqlName = strings.Replace(graphqlName, "URL", "Url", -1)

	return strcase.ToLowerCamel(graphqlName)
}

func toGraphQLType(fieldName, boilerType string) string {
	lowerFieldName := strings.ToLower(fieldName)
	lowerBoilerType := strings.ToLower(boilerType)
	if strings.Contains(lowerBoilerType, "string") {
		return "String"
	}
	if strings.Contains(lowerBoilerType, "int") {
		if strings.HasSuffix(lowerFieldName, "id") {
			return "ID"
		}
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

func isFirstCharacterLowerCase(s string) bool {
	if len(s) > 0 && s[0] == strings.ToLower(s)[0] {
		return true
	}
	return false
}

func appendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}
