package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/iancoleman/strcase"
	"github.com/micro/cli"
	"github.com/web-ridge/gqlgen-sqlboiler/boiler"
)

type Model struct {
	Name       string
	Implements *string
}
type Field struct {
	Name     string
	Required bool
	Type     string // String, ID, Integer
}

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
			fmt.Println(outputFile)
			// boilerTypeMap, _ := boiler.ParseBoilerFile(modelDirectory)

			// boilerTypeMap =
			// Block.L blockL
			// FlowBlock.ID int
			// User.LastName null.String

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
				fieldName := splitted[1]
				if isFirstCharacterLowerCase(modelName) {

					// It's the relations of the model
					// let's add them so we can use them later
					
					if strings.HasSuffix(modelName, "R") {
						modelName = strcase.ToCamel(strings.TrimSuffix(modelName, "R"))
						_, ok := relationsPerModel[modelName]
						if !ok {
							relationsPerModel[modelName] = []*Field{}
						}
						fmt.Println("adding relation "+fieldName+" to " + modelName + " ")
						relationsPerModel[modelName] = append(relationsPerModel[modelName], &Field{
							Name: fieldName,
							Type: boilerType,
						})
					} 


					continue

				}

				_, ok := fieldPerModel[modelName]
				if !ok {
					fieldPerModel[modelName] = []*Field{}
				}

				if fieldName == "L" || fieldName == "R" {
					continue
				}

				fieldPerModel[modelName] = append(fieldPerModel[modelName], &Field{
					Name: fieldName,
					Type: boilerType,
				})
			}

			for modelName, relations := range relationsPerModel {
				for _, relationField := range relations {
					// remove relationID from fields since the relationship is preloadable
					fieldPerModel[modelName] = append(fieldPerModel[modelName], relationField)

				}
			}

			for model, fields := range fieldPerModel {
				fmt.Println("type " + model + " {")
				for _, field := range fields {
					fmt.Println("        " + field.Name + " : " + field.Type)
				}
				fmt.Println("}")
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func isFirstCharacterLowerCase(s string) bool {
	if len(s) > 0 && s[0] == strings.ToLower(s)[0] {
		return true
	}
	return false
}
