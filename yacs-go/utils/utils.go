package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var hasDigitalOnly *regexp.Regexp

func init() {
	hasDigitalOnly = regexp.MustCompile(`^[0-9]+$`)
}

// FileForProcess is object for stroing information processing of single file
type FileForProcess struct {
	From, To string
	Short    string
	Num      int
}

// FindAllFiles returns list of files for processing
func FindAllFiles(dirFrom, dirTo, subDirToSkip string) ([]FileForProcess, error) {

	res := []FileForProcess{}

	err := filepath.Walk(dirFrom, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return fmt.Errorf("prevent panic by handling failure accessing a path %q: %v", dirFrom, err)
		}

		if info.IsDir() && info.Name() == subDirToSkip {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		// fmt.Printf("visitedfile: %q\n", path)

		short := strings.TrimPrefix(path, dirFrom)

		n := FileForProcess{
			From:  path,
			To:    filepath.Join(dirTo, short),
			Short: short,
		}

		res = append(res, n)

		return nil
	})

	return res, err
}

// CreateDirIfNotExist create dir if it doesn't exist or return error.
func CreateDirIfNotExist(dir string, mode os.FileMode) error {
	/* CreateDirIfNotExist("/User/Ivan/Susanin", os.FileMode(0666)) */
	fi, err := os.Stat(dir)

	if os.IsPermission(err) {
		return err
	}

	if os.IsNotExist(err) {
		// os.Mkdir(dir, os.FileMode(0666))
		return os.MkdirAll(dir, mode)
	}

	if strings.TrimLeft(fi.Mode().String(), "-d") != strings.TrimLeft(mode.String(), "-d") {
		return os.Chmod(dir, mode)
	}

	return nil
}

// URLDefrag splits uri to fragment.
func URLDefrag(URI string) (string, string, error) {
	u, err := url.Parse(URI)
	if err != nil {
		return "", "", err
	}

	fragment := u.Fragment
	u.Scheme = ""
	u.Fragment = ""
	return filepath.Clean(u.String()), fragment, nil
}

// SaveJSONFile stores interface to json file.
func SaveJSONFile(file string, data interface{}, mode os.FileMode) error {

	dir, _ := filepath.Split(file)
	err := CreateDirIfNotExist(dir, mode)
	if err != nil {
		return err
	}

	json, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	json = append(json, ([]byte("\n"))...)

	err = ioutil.WriteFile(file, json, mode)
	return err
}

// GetKeyFromIntefaceString returns string value from hash.
func GetKeyFromIntefaceString(node interface{}, key string) (string, bool) {
	if i, find := GetKeyFromInteface(node, key); find {
		switch i.(type) {
		case string:
			return i.(string), true
		}
	}
	return "", false
}

// DoesIntefaceHaveKey checks key in hash.
func DoesIntefaceHaveKey(node interface{}, key string) bool {
	find := false

	switch node.(type) {
	case map[string]interface{}:
		_, find = (node.(map[string]interface{}))[key]
		return find
	case map[string]string:
		_, find = (node.(map[string]string))[key]
		return find
	}
	return find
}

// GetKeyFromAnyInteface returns value from hash by key and from array by index.
func GetKeyFromAnyInteface(node interface{}, key string) (interface{}, bool) {

	if node == nil {
		return nil, false
	}

	switch node.(type) {

	case map[string]interface{}:
		if res, find := (node.(map[string]interface{}))[key]; find {
			return res, find
		}

	case []interface{}:

		iStr := hasDigitalOnly.FindString(key)
		if iStr == "" {
			return nil, false
		}

		i, err := strconv.Atoi(iStr)
		if err != nil || i < 0 {
			return nil, false
		}

		m := node.([]interface{})
		if len(m) <= i {
			return nil, false
		}

		return m[i], true
	}

	return nil, false
}

// GetKeyFromInteface returns value from hash by key.
func GetKeyFromInteface(node interface{}, key string) (interface{}, bool) {

	if node == nil {
		return nil, false
	}

	switch node.(type) {
	case map[string]interface{}:
		if res, find := (node.(map[string]interface{}))[key]; find {
			return res, find
		}
	}
	return nil, false
}

// CheckBoolInteface checks any key from array in hash.
func CheckBoolInteface(node interface{}, keys []string, toCheck bool) bool {

	if len(keys) == 0 {
		return false
	}

	find := false
	for _, key := range keys {
		node, find = GetKeyFromInteface(node, key)
	}

	if find {
		switch node.(type) {
		case bool:
			return toCheck == node.(bool)
		}
	}

	return false
}

// IsListInterface checks tyep as []interface{}
func IsListInterface(doc interface{}) bool {
	switch doc.(type) {
	case []interface{}, []string:
		return true
	}
	return false
}

// IsMapStringInterface checks tyep as map[string]interface{} or map[string]string
func IsMapStringInterface(doc interface{}) bool {

	switch doc.(type) {
	case map[string]interface{}, map[string]string:
		return true
	}
	return false
}

// ToListInterface convert interface{} to []interface{}
func ToListInterface(doc interface{}) []interface{} {

	switch doc.(type) {

	case []interface{}:
		return doc.([]interface{})
	case []string:
		out := []interface{}{}
		for k, v := range doc.([]string) {
			out[k] = v
		}
		return out
	}

	return []interface{}{doc}
}

// ReversedInterface reverses []interface{}
func ReversedInterface(doc interface{}) []interface{} {

	switch doc.(type) {
	case []interface{}:
		m := doc.([]interface{})
		for i, j := 0, len(m)-1; i <= j; i, j = i+1, j-1 {
			m[i], m[j] = m[j], m[i]
		}
		return m
	}

	return []interface{}{doc}
}

// DeepCopy copies interface{}
func DeepCopy(node interface{}) (interface{}, error) {

	if node == nil {
		return nil, nil
	}

	var err error

	switch node.(type) {

	case bool, int64, float64, string:
		return node, nil

	case []bool:
		return node.([]bool)[:], nil

	case []int64:
		return node.([]int64)[:], nil

	case []float64:
		return node.([]float64)[:], nil

	case []string:
		return node.([]string)[:], nil

	case map[string]interface{}:

		out := map[string]interface{}{}
		for k, v := range node.(map[string]interface{}) {
			out[k], err = DeepCopy(v)
			if err != nil {
				return nil, err
			}
		}
		return out, nil

	case []interface{}:
		out := []interface{}{}
		for _, v := range node.([]interface{}) {
			next, err := DeepCopy(v)
			if err != nil {
				return nil, err
			}
			out = append(out, next)
		}
		return out, nil
	}

	return node, fmt.Errorf("DeepCopy(...): bad type of interface - %s", reflect.TypeOf(node))
}
