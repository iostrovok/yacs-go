package helper

/*

Implements the JSON Reference IETF Draft: https://tools.ietf.org/html/draft-pbryan-zyp-json-ref-03
A JSON Reference is used to insert a value from another JSON document as the
value for this key.  For instance if this document had its JSON References
resolved:
    {
        "gannett_founder": "Frank Gannett",
        "hearst_rival": {"$ref": "///gannett_founder"}
    }
it would look like:
    {
        "gannett_founder": "Frank Gannett",
        "hearst_rival": "Frank Gannett"
    }
JSON References can also refer to other documents by giving a full URI such as:
    {
        "hearst_rival": {"$ref": "http://example.com/gannett.json///founder"}
    }

*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"myconst"
	"utils"
)

func isJSONRef(doc interface{}) bool {
	// Tests to see if this is a JSON Reference
	return utils.DoesIntefaceHaveKey(doc, myconst.JSONRefKeyName)
}

func extractJSONrefURI(value interface{}) (string, error) {
	// Returns the JSON Reference URI from a JSON Reference //
	out, find := utils.GetKeyFromIntefaceString(value, myconst.JSONRefKeyName)
	if !find {
		return "", fmt.Errorf("'%s' is not found", myconst.JSONRefKeyName)
	}
	return out, nil
}

func fetchURL(url string) ([]byte, error) {
	// response = urllib2.urlopen(url)
	// return response.read()
	return ioutil.ReadFile(url)
}

func fetchURI(uri string, context *Context) (interface{}, error) {
	// Fetch URI as JSON.
	// url is what we'll actually end up retrieving
	// url := filepath.Clean(filepath.Join(context.getDir(), uri))
	url := context.getDir(uri)

	uriData, err := fetchURL(url)
	if err != nil {
		return nil, err
	}

	var doc interface{}
	err = json.Unmarshal(uriData, &doc)
	if err == nil {
		// We just retrieved a new URL so the context has changed.
		context.setURI(url)
	}

	return doc, err
}

func getRefURI(uri string, docIn interface{}, context *Context) (interface{}, error) {
	/* Returns the value designated by the provided JSON Reference URI.

	   If the reference is local, retrieves it from 'context' parameter.
	   If the reference is not local, consults the registry for the document.
	*/

	doc, err := utils.DeepCopy(docIn)
	if err != nil {
		return nil, err
	}

	// JSON Reference URIs look mostly like URLs.
	base, pointer, err := utils.URLDefrag(uri)
	if err != nil {
		return nil, err
	}

	// There's a path, so we need to fetch the document.
	if base != "" {
		doc, err = fetchURI(base, context)
		if err != nil {
			return nil, err
		}
	}

	// The pointer piece of a JSON Reference appears after an (optional) hash.
	if pointer == "" {
		return doc, nil
	}

	return resolvePointer(doc, pointer)
}

func resolvePointer(node interface{}, path string) (interface{}, error) {

	paths := strings.Split(strings.TrimPrefix(strings.TrimSuffix(path, "/"), "/"), "/")
	if path == "/" || path == "" || len(paths) == 0 {
		return utils.DeepCopy(node)
	}

	next, find := utils.GetKeyFromAnyInteface(node, paths[0])
	if !find {
		return nil, nil
	}

	for _, v := range paths[1:] {
		next, find = utils.GetKeyFromAnyInteface(next, v)
		if !find {
			return nil, nil
		}
	}

	return utils.DeepCopy(next)
}
