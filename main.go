package main

import (
	"fmt"
	"log"
	"os"
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

type Model struct {
	Name       string
	Implements *string
	Fields     []Field
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
			schema := getSchema(modelDirectory, mutations, batchUpdate, batchCreate, batchDelete)
			fmt.Println(schema)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

const queryHelperStructs = `
type FilterID {
	in: [ID!]
	notIn: [ID!]
}
type FilterString {
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
type FilterInteger {
	equalTo: Int
	lessThan: Int
	lessThanOrEqualTo: Int
	moreThan: Int
	moreThanOrEqualTo: Int
	in: [Int!]
	notIn: [Int!]
}
type FilterBoolean {
	equalTo: Boolean
}
`

func getSchema(
	modelDirectory string,
	mutations bool,
	batchUpdate bool,
	batchCreate bool,
	batchDelete bool,
) string {
	var schema strings.Builder

	boilerTypeMap, _, boilerTypeOrder := boiler.ParseBoilerFile(modelDirectory)

	fieldPerModel := make(map[string][]*Field)
	relationsPerModel := make(map[string][]*Field)

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
		boilerType := boilerTypeMap[modelAndField]
		splitted := strings.Split(modelAndField, ".")
		modelName := splitted[0]
		boilerFieldName := splitted[1]
		if isFirstCharacterLowerCase(modelName) {

			// It's the relations of the model
			// let's add them so we can use them later

			if strings.HasSuffix(modelName, "R") {
				modelName = strcase.ToCamel(strings.TrimSuffix(modelName, "R"))
				_, ok := relationsPerModel[modelName]
				if !ok {
					relationsPerModel[modelName] = []*Field{}
				}
				// fmt.Println("adding relation " + fieldName + " to " + modelName + " ")
				relationField := toField(boilerFieldName, boilerType)
				relationField.IsRelation = true
				relationsPerModel[modelName] = append(relationsPerModel[modelName], relationField)
			}

			continue

		}

		_, ok := fieldPerModel[modelName]
		if !ok {
			fieldPerModel[modelName] = []*Field{}
		}

		if boilerFieldName == "L" || boilerFieldName == "R" {
			continue
		}

		fieldPerModel[modelName] = append(fieldPerModel[modelName], toField(boilerFieldName, boilerType))
	}

	// relations
	for modelName, relations := range relationsPerModel {
		fieldPerModel[modelName] = append(fieldPerModel[modelName], relations...)
	}

	// take care that models are always returned in same order
	sortedModelNames := make([]string, 0, len(relationsPerModel))
	for k := range relationsPerModel {
		sortedModelNames = append(sortedModelNames, k)
	}
	sort.Strings(sortedModelNames)

	for _, model := range sortedModelNames {
		fields := relationsPerModel[model]

		schema.WriteString("type " + model + " {")
		schema.WriteString("\n")
		for _, field := range fields {

			gType := field.Type

			if field.IsArray {
				// To use a list type, surround the type in square brackets, so [Int] is a list of integers.
				gType = "[" + gType + "]"
			}
			if field.IsRequired {
				// Use an exclamation point to indicate a type cannot be nullable,
				// so String! is a non-nullable string.
				gType = gType + "!"
			}
			field.FullType = gType

			schema.WriteString(indent + field.Name + " : " + gType)
			schema.WriteString("\n")
		}
		schema.WriteString("}")
		schema.WriteString("\n")
	}

	// Add helpers for filtering lists
	schema.WriteString(queryHelperStructs)
	schema.WriteString("\n")

	// generate filter structs per model
	for _, model := range sortedModelNames {
		fields := relationsPerModel[model]

		// Generate a type safe grapql filter

		// Generate the base filter
		// type UserFilter {
		// 	search: String
		// 	where: UserWhere
		// }
		schema.WriteString("type " + model + "Filter {")
		schema.WriteString("\n")
		schema.WriteString(indent + "search: String")
		schema.WriteString("\n")
		schema.WriteString(indent + "where: " + model + "Where")
		schema.WriteString("\n")
		schema.WriteString("}")
		schema.WriteString("\n")
		// Generate a where struct
		// type UserWhere {
		// 	id: IDFilter
		// 	title: StringFilter
		// 	organization: OrganizationWhere
		// 	or: FlowBlockWhere
		// 	and: FlowBlockWhere
		// }
		schema.WriteString("type " + model + "Where {")
		schema.WriteString("\n")
		for _, field := range fields {
			if field.IsRelation {
				// Support filtering in relationships (atleast schema wise)
				schema.WriteString(indent + field.Name + ": " + field.Type + "Where")
				schema.WriteString("\n")
			} else {
				schema.WriteString(indent + field.Name + ": " + field.Type + "Filter")
				schema.WriteString("\n")
			}

		}
		schema.WriteString(indent + "or: " + model + "Where")
		schema.WriteString("\n")

		schema.WriteString(indent + "and: " + model + "Where")
		schema.WriteString("\n")

		schema.WriteString("}")
		schema.WriteString("\n")
	}

	schema.WriteString("type Query {")
	schema.WriteString("\n")
	for _, model := range sortedModelNames {
		// single models
		schema.WriteString(indent)
		schema.WriteString(strcase.ToLowerCamel(model) + "(id: ID!)")
		schema.WriteString(": ")
		schema.WriteString(model + "!")
		schema.WriteString("\n")

		// lists
		modelArray := pluralizer.Plural(model)
		schema.WriteString(indent)
		// TODO: pagination
		schema.WriteString(strcase.ToLowerCamel(modelArray) + "(filter: " + model + "Filter)")
		schema.WriteString(": ")
		schema.WriteString("[" + model + "!]!")
		schema.WriteString("\n")

	}
	schema.WriteString("}")
	schema.WriteString("\n")

	// Generate input and payloads for mutatations
	if mutations {
		for _, model := range sortedModelNames {
			fields := relationsPerModel[model]
			modelArray := pluralizer.Plural(model)
			// input UserInput {
			// 	firstName: String!
			// 	lastName: String
			//	organizationId: ID!
			// }
			schema.WriteString("input " + model + "Input {")
			schema.WriteString("\n")
			for _, field := range fields {
				// not possible yet in input
				if field.IsRelation {
					continue
				}
				schema.WriteString(indent + field.Name + ": " + field.FullType)
				schema.WriteString("\n")
			}
			schema.WriteString("}")
			schema.WriteString("\n")

			// type UserPayload {
			// 	user: User!
			// }
			schema.WriteString("type " + model + "Payload {")
			schema.WriteString("\n")
			schema.WriteString(indent + strcase.ToLowerCamel(model) + ": " + model + "!")
			schema.WriteString("\n")
			schema.WriteString("}")
			schema.WriteString("\n")

			// TODO batch, delete input and payloads

			// type UserDeletePayload {
			// 	id: ID!
			// }
			schema.WriteString("type " + model + "DeletePayload {")
			schema.WriteString("\n")
			schema.WriteString(indent + "id: ID!")
			schema.WriteString("\n")
			schema.WriteString("}")
			schema.WriteString("\n")
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
			}

		}
	}

	// Generate mutation queries

	if mutations {
		schema.WriteString("type Mutation {")
		schema.WriteString("\n")
		for _, model := range sortedModelNames {

			modelArray := pluralizer.Plural(model)

			// create single
			// e.g createUser(input: UserInput!): UserPayload!
			schema.WriteString(indent)
			schema.WriteString("create" + model + "(input: " + model + "Input!)")
			schema.WriteString(": ")
			schema.WriteString(model + "Payload!")
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
			schema.WriteString("update" + model + "(input: " + model + "Input!)")
			schema.WriteString(": ")
			schema.WriteString(model + "Payload!")
			schema.WriteString("\n")

			// update multiple (batch update)
			// e.g updateUsers(filter: UserFilter, input: [UsersInput!]!): UsersPayload!
			if batchUpdate {
				schema.WriteString(indent)
				schema.WriteString("update" + modelArray + "(filter: " + model + "Filter, input: " + modelArray + "Input!)")
				schema.WriteString(": ")
				schema.WriteString(modelArray + "UpdatePayload!")
				schema.WriteString("\n")
			}

			// delete single
			// e.g deleteUser(id: ID!): UserPayload!
			schema.WriteString(indent)
			schema.WriteString("delete" + model + "(id: ID!)")
			schema.WriteString(": ")
			schema.WriteString(model + "DeletePayload!")
			schema.WriteString("\n")

			// delete multiple
			// e.g deleteUsers(filter: UserFilter, input: [UsersInput!]!): UsersPayload!
			if batchDelete {
				schema.WriteString(indent)
				schema.WriteString("delete" + modelArray + "(filter: " + model + "Filter)")
				schema.WriteString(": ")
				schema.WriteString(modelArray + "DeletePayload!")
				schema.WriteString("\n")
			}

		}
		schema.WriteString("}")
		schema.WriteString("\n")
	}

	return schema.String()
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
func toField(boilerName, boilerType string) *Field {
	return &Field{
		Name:       toGraphQLName(boilerName, boilerType),
		Type:       toGraphQLType(boilerName, boilerType),
		BoilerName: boilerName,
		BoilerType: boilerType,
		IsRequired: isRequired(boilerType),
		IsArray:    isArray(boilerType),
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
