package helper

import (
	"jsonschema"
	"myconst"
	"utils"
)

func notNeedResolv(doc interface{}) bool {
	/*
	   check
	   ------
	      "@doc": {
	          "resolve": false
	      }
	   ------
	*/
	// Don't process this doc if resolve is turned off via directive
	if utils.DoesIntefaceHaveKey(doc, myconst.DocKeyName) {
		if utils.CheckBoolInteface(doc, []string{myconst.DocKeyName, myconst.ResolveKeyName}, false) {
			return true
		}
	}
	return false
}

func resolveDoc(doc interface{}, context *Context) (interface{}, error) {

	// Don't process this doc if resolve is turned off via directive
	if notNeedResolv(doc) {
		return doc, nil
	}

	// Now deal with potentially failing reference retrievals in dicts/lists
	if utils.IsMapStringInterface(doc) {
		if isJSONRef(doc) {
			// Look for JSON References in all pieces of this document.
			return resolveRefInDoc(doc, context)
		}
		return resolveMapDoc(doc.(map[string]interface{}), context)
	}

	if utils.IsListInterface(doc) {
		return resolveArrayDoc(doc.([]interface{}), context)
	}

	return doc, nil
}

func resolveRefInDoc(docIn interface{}, context *Context) (interface{}, error) {

	uri, err := extractJSONrefURI(docIn)
	if err != nil {
		return nil, err
	}

	absuri := context.getDir(uri)
	if out, find := uriCache.Get(absuri); find {
		return out, nil
	}

	// Context will be getting updated so make a copy.
	docContext := context.copy()
	doc, err := getRefURI(uri, docIn, docContext)
	if err != nil {
		return nil, err
	}

	out, err := resolveDoc(doc, docContext)
	if err != nil {
		return nil, err
	}

	uriCache.Add(absuri, out)
	return out, nil
}

func resolveMapDoc(doc map[string]interface{}, context *Context) (interface{}, error) {
	out := map[string]interface{}{}
	for key, value := range doc {
		res, err := resolveDoc(value, context)
		if err != nil {
			return nil, err
		}
		out[key] = res
	}
	return out, nil
}

func resolveArrayDoc(doc []interface{}, context *Context) (interface{}, error) {
	out := []interface{}{}
	for _, value := range doc {
		res, err := resolveDoc(value, context)
		if err != nil {
			return nil, err
		}
		out = append(out, res)
	}
	return out, nil
}

func processDoc(doc interface{}, context *Context, needToResolve, needToInherit, validateJSONSchema, verbose bool) (interface{}, error) {

	var err error
	processed := doc

	// Resolve References
	if needToResolve {
		processed, err = resolveDoc(processed, context)
		if err != nil {
			return nil, err
		}
	}

	// Apply Inheritance/locking
	if needToInherit {
		processed, err = mergeParents(processed, true)
		if err != nil {
			return nil, err
		}
	}

	// Validate schema if possible
	if !validateJSONSchema {
		return jsonschema.RemoveSchemaReferences(processed), nil
	}

	return jsonschema.ValidateSchema(processed, verbose)
}

// Process is main function - it starts processing single files by URI
func Process(uri string, needToResolve, needToInherit, validateJSONSchema, verbose bool) (interface{}, error) {

	context := newContext()

	doc, err := getRefURI(uri, nil, context)
	if err != nil {
		return nil, err
	}

	// Resolve references + inherit
	return processDoc(doc, context, needToResolve, needToInherit, validateJSONSchema, verbose)
}
