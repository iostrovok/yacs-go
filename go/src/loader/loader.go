package loader

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// GetURI returns json parsed object by URL or file link.
func GetURI(filename string) (interface{}, error) {

	var body []byte
	var err error

	if isURL(filename) {
		body, err = getURL(filename)
	} else {
		body, err = loadFile(filename)
	}

	if err != nil {
		return nil, err
	}

	var s interface{}
	err = json.Unmarshal(body, &s)
	return s, err
}

func loadFile(filename string) ([]byte, error) {

	filename = strings.TrimPrefix(filename, "file://")

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(file)
}

func getURL(filename string) ([]byte, error) {

	resp, err := http.Get(filename)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// isValidUrl tests a string to determine if it is a url or not.
func isURL(testURL string) bool {
	//_, err := url.ParseRequestURI(toTest)
	u, err := url.Parse(testURL)

	if err != nil {
		return false
	}

	return u.Scheme == "http" || u.Scheme == "https"
}
