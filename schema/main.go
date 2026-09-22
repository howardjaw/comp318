package main

import (
	"fmt"
	"os"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: go run . <data.json> <schema.json>")
		return
	}
	dataFile := os.Args[1]
	schemaFile := os.Args[2]

	// Load and compile schema
	compiler := jsonschema.NewCompiler()

	schema, err := compiler.Compile(schemaFile)
	if err != nil {
		fmt.Println("cannot load or compile json schema: ", err)
		return
	}

	// Open the data file
	file, err := os.Open(dataFile)
	if err != nil {
		fmt.Println("Could not open data file: ", err)
		return
	}
	defer file.Close()

	// Unmarshal JSON
	data, err := jsonschema.UnmarshalJSON(file)
	if err != nil {
		fmt.Println("Could not unmarshal JSON: ", err)
		return
	}

	// Validate JSON schema
	err = schema.Validate(data)
	if err != nil {
		fmt.Println("Invalid JSON schema: ", err)
		return
	}

	fmt.Println("Valid JSON schema")
}
