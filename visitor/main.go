package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	jsonvalue "github.com/howardjaw/comp318/jsonvalue"
)

/**
Implementation of the visitor
**/

type StringVisitor struct{} //An empty struct

// Give the empty struct methods

func (v StringVisitor) Bool(value bool) (string, error) {
	return strconv.FormatBool(value), nil
}

func (v StringVisitor) Float64(value float64) (string, error) {
	return strconv.FormatFloat(value, 'g', -1, 64), nil
}

func (v StringVisitor) Null() (string, error) {
	return "null", nil
}

func (v StringVisitor) String(value string) (string, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

func (v StringVisitor) Slice(values []jsonvalue.JSONValue) (string, error) {
	parts := make([]string, 0, len(values))

	for _, child := range values {
		text, err := jsonvalue.Accept[string](child, v)
		if err != nil {
			return "", err
		}
		parts = append(parts, text)
	}
	return "[" + strings.Join(parts, ",") + "]", nil
}

func (v StringVisitor) Map(values map[string]jsonvalue.JSONValue) (string, error) {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	// sort.Strings(keys)

	parts := make([]string, 0, len(values))

	for _, key := range keys {
		keyText, err := v.String(key)
		if err != nil {
			return "", err
		}

		valueText, err := jsonvalue.Accept[string](values[key], v)
		if err != nil {
			return "", err
		}

		parts = append(parts, keyText+": "+valueText)
	}

	return "{" + strings.Join(parts, ",") + "}", nil
}

/**
Main function calls visitor
**/

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run . <data.json>")
		return
	}

	contents, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Could not read file: ", err)
	}

	var data jsonvalue.JSONValue

	err = json.Unmarshal(contents, &data)
	if err != nil {
		fmt.Println("Could not parse JSON: ", err)
		return
	}

	visitor := StringVisitor{}

	result, err := jsonvalue.Accept[string](data, visitor)
	if err != nil {
		fmt.Println("Could not process json: ", err)
		return
	}

	fmt.Println(result)
}
