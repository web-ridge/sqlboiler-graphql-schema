package main

import (
	"fmt"
	pluralize "github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/micro/cli"
	"github.com/web-ridge/gqlgen-sqlboiler/boiler"
	"log"
	"os"
	"strings"
)

var pluralizer *pluralize.Client

func init() {
	pluralizer = pluralize.NewClient()
}

const indent = "  "

type Model struct {
	Name       string
	Implements *string
	Fields     []Field
}
type Field struct {
	Name       string
	Type       string // String, ID, Integer
	BoilerName string
	BoilerType string
	IsRequired bool
	IsArray    bool
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

func main() {
	var modelDirectory string
	var outputFile string

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
		},
		Action: func(c *cli.Context) error {
			var schema strings.Builder

			boilerTypeMap, _ := boiler.ParseBoilerFile(modelDirectory)
			// fmt.Println("boilerStructMap")
			// for name, value := range boilerStructMap {
			// 	fmt.Println(name, value)

			// }
			// fmt.Println("")
			// fmt.Println("")
			// fmt.Println("")
			// fmt.Println("")
			// fmt.Println("boilerTypeMap")
			// for name, value := range boilerTypeMap {
			// 	fmt.Println(name, value)

			// }

			fieldPerModel := make(map[string][]*Field)
			relationsPerModel := make(map[string][]*Field)

			for modelAndField, boilerType := range boilerTypeMap {
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
						relationsPerModel[modelName] = append(relationsPerModel[modelName], toField(boilerFieldName, boilerType))
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

			for modelName, relations := range relationsPerModel {
				fieldPerModel[modelName] = append(fieldPerModel[modelName], relations...)
			}
			for model, fields := range fieldPerModel {

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
			for model, fields := range fieldPerModel {

				// Generate a type safe grapql filter

				// Generate the base filter
				// type UserFilter {
				// 	search: String
				// 	where: UserWhere
				// }

				// Generate a where struct
				// type UserWhere {
				// 	id: IDFilter
				// 	title: StringFilter
				// 	or: FlowBlockWhere
				// 	and: FlowBlockWhere
				// }
				schema.WriteString("type " + model + "Where {")
				schema.WriteString("\n")
				for _, field := range fields {
					schema.WriteString(indent + field.Name + ": " + field.Type + "Filter")
					schema.WriteString("\n")
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
			for model, _ := range fieldPerModel {

				// single models
				schema.WriteString(indent)
				schema.WriteString(strcase.ToLowerCamel(model) + "(id: ID!)")
				schema.WriteString(":")
				schema.WriteString(model + "!")
				schema.WriteString("\n")

				// lists
				modelArray := pluralizer.Plural(model)
				schema.WriteString(indent)
				schema.WriteString(strcase.ToLowerCamel(modelArray) + "")
				schema.WriteString(":")
				schema.WriteString("[" + model + "!]!")
				schema.WriteString("\n")

			}
			schema.WriteString("}")
			schema.WriteString("\n")
			fmt.Println(schema.String())

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

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

	// e.g. OrganizationID
	graphqlName = strings.Replace(graphqlName, "ID", "Id", -1)

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
