package main

import (
	"fmt"

	"github.com/web-ridge/gqlgen-sqlboiler/boiler"
)

func main() {
	var backendDir string

	boilerTypeMap, boilerStructMap := boiler.ParseBoilerFile(backendDir)
	fmt.Println(boilerTypeMap)
	fmt.Println(boilerStructMap)

}
