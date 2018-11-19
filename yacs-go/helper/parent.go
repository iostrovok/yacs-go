package helper

import (
	"myconst"
	"utils"
)

// getLockNames extracts string keys from interface and convert to map[string]bool
func getLockNames(doc interface{}) map[string]bool {
	/*
		It ignores bool, int64, float64, []bool, []int64, []float64.

		Exmaples:

			"@lock_names" : "header"

			"@lock_names" : [
				"header",
				"footer"
			]

			"@lock_names" : {
				"header": "1",
				"footer": "1"
			}


			"@lock_names" : {
				"header": {},
				"footer": true
			}
	*/

	lockNames, find := utils.GetKeyFromInteface(doc, myconst.LockKeyName)
	if !find {
		return map[string]bool{}
	}

	out := map[string]bool{}
	switch lockNames.(type) {

	case bool, int64, float64, []bool, []int64, []float64:
		return out

	case string:
		out[lockNames.(string)] = true

	case []string:
		for _, v := range lockNames.([]string) {
			out[v] = true
		}

	case map[string]interface{}:
		for k := range lockNames.(map[string]interface{}) {
			out[k] = true
		}

	case []interface{}:
		for _, v := range lockNames.([]interface{}) {
			switch v.(type) {
			case string:
				out[v.(string)] = true
			}
		}
	}

	return out
}

func replaceExceptLocked(original, replacement interface{}) interface{} {
	/*
	   Overwrites values in the original dict with those in the replacement dict,
	   unless the key is in @lock_names.
	*/

	if utils.IsMapStringInterface(original) && utils.IsMapStringInterface(replacement) {

		// See comment to getLockNames() function
		lockNames := getLockNames(original)

		out := original.(map[string]interface{})

		// Iterate through the replacement dictionary.
		for key, value := range replacement.(map[string]interface{}) {
			if lockNames[key] {
				/*
					Don't process key from "@lock_names" keys.
					See comment to getLockNames() function
				*/
				continue
			}

			originalValues, originalFind := utils.GetKeyFromInteface(original, key)

			if myconst.SchemaKeyName == key {
				if value != nil {
					out[key] = value
				} else {
					out[key] = originalValues
				}
				continue
			}

			// If the values are a dictionary...
			if originalFind && originalValues != nil && utils.IsMapStringInterface(value) {
				// recursive call is processing viscera of structures like map[string]interface{}
				out[key] = replaceExceptLocked(originalValues, value)
			} else {
				// save the value to the final response.
				out[key] = value
			}
		}

		return out
	}

	return replacement
}

func mergeParents(doc interface{}, removeParentRef bool) (interface{}, error) {
	// Merges any included parent documents with this document.

	// Only actually try to merge if this is a dict
	switch doc.(type) {
	case map[string]interface{}:

		m := doc.(map[string]interface{})

		// Merge with parent(s) at this level
		if parentsIn, find := m[myconst.ParentKeyName]; find {
			delete(m, myconst.ParentKeyName)

			// Can be a single parent or list of them, so normalize to list
			parents := utils.ToListInterface(parentsIn)

			// Evaluate from left to right: first in list has last priority
			next, err := mergeParents(parents, removeParentRef)
			if err != nil {
				return nil, err
			}

			doc = m
			for _, parent := range utils.ReversedInterface(next) {
				doc = replaceExceptLocked(parent, doc)
			}
		}

		switch doc.(type) {
		case map[string]interface{}:
			m := doc.(map[string]interface{})
			// Recursively check any children and merge them too
			for key, value := range m {
				v, err := mergeParents(value, removeParentRef)
				if err != nil {
					return nil, err
				}

				m[key] = v
			}

			return m, nil
		}

		return doc, nil

	// Process items in a list
	case []interface{}:

		out := []interface{}{}
		for _, listItem := range doc.([]interface{}) {
			res, err := mergeParents(listItem, removeParentRef)
			if err != nil {
				return nil, err
			}
			out = append(out, res)
		}

		return out, nil
	}

	// If not a dict or a list, just the value
	return doc, nil

}
