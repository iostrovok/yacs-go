package jsonschema

/*
See draft:
https://tools.ietf.org/html/draft-pbryan-zyp-json-ref-03

Example usage:

/config-files/my.json
    {
        "@schemas": {
            "anykey" : {
                "$ref": "../../schemas/schemas-file.json"
            }
        },
        "username": "web"
    }

/schemas/schemas-file.json
    {
        "$schema": "http://json-schema.org/draft-04/schema#",
        "required" : [
            "username"
        ],
        "type" : "object",
        "properties" : {
            "username" : {
                "type" : "string"
            }
        }
    }

*/

import (
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"

	"github.com/iostrovok/yacs-go/yacs-go/myconst"
	"github.com/iostrovok/yacs-go/yacs-go/utils"
)

func joinErrors(errors []error) error {
	if len(errors) == 0 {
		return nil
	}

	strs := []string{}
	for i, e := range errors {
		strs = append(strs, fmt.Sprintf("[%d] %s", i, e))
	}

	return fmt.Errorf(strings.Join(strs, "\n"))
}

// ValidateSchema validates json schema
func ValidateSchema(doc interface{}, verbose bool) (interface{}, error) {
	out, errs := validateListSchemas(doc, verbose)
	return out, joinErrors(errs)
}

func validateListSchemas(doc interface{}, verbose bool) (interface{}, []error) {

	var e []error

	switch doc.(type) {
	case map[string]interface{}:
		out, errors := validateAgainstSchemas(doc.(map[string]interface{}), verbose)
		for key, value := range doc.(map[string]interface{}) {
			out[key], e = validateListSchemas(value, verbose)
			errors = append(errors, e...)
		}
		return out, errors

	case []interface{}:
		out := make([]interface{}, len(doc.([]interface{})))
		errors := []error{}
		for i, value := range doc.([]interface{}) {
			out[i], e = validateListSchemas(value, verbose)
			errors = append(errors, e...)
		}
		return out, errors
	}
	return doc, []error{}
}

func validateAgainstSchemas(doc map[string]interface{}, verbose bool) (map[string]interface{}, []error) {

	out := []error{}

	schemas, find := utils.GetKeyFromInteface(doc, myconst.SchemaKeyName)
	if !find {
		return doc, out
	}

	delete(doc, myconst.SchemaKeyName)

	switch schemas.(type) {
	case map[string]interface{}:
		for _, schema := range schemas.(map[string]interface{}) {
			switch schema.(type) {
			case map[string]interface{}:
				validateOneSchema(schema.(map[string]interface{}), doc, verbose)
			}
		}
	}

	return doc, out
}

func validateOneSchema(schema, body map[string]interface{}, verbose bool) error {

	schemaLoader := gojsonschema.NewGoLoader(schema)
	documentLoader := gojsonschema.NewGoLoader(body)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}

	if result.Valid() {
		if verbose {
			fmt.Printf("The document is valid\n")
		}
		return nil
	}

	out := "The document is not valid. see errors :\n"
	out += fmt.Sprintf("schema: %v\n", schema)

	for _, desc := range result.Errors() {
		out += fmt.Sprintf("- %s\n", desc)
	}

	return fmt.Errorf(out)
}

// RemoveSchemaReferences - removes "@schemas" objects from JSON without validation.
func RemoveSchemaReferences(doc interface{}) interface{} {

	if doc == nil {
		return nil
	}

	switch doc.(type) {
	case map[string]interface{}:

		out := doc.(map[string]interface{})
		delete(out, myconst.SchemaKeyName)
		for key, value := range out {
			out[key] = RemoveSchemaReferences(value)
		}
		return out

	case []interface{}:
		out := doc.([]interface{})
		for i := range out {
			out[i] = RemoveSchemaReferences(out[i])
		}
		return out
	}
	return doc
}
