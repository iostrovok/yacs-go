package diff

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
)

// Diff is interface function for checking difference between 2 JSON objects
func Diff(a, b interface{}) []string {
	res := deepDiff(a, b, "")
	return sortAndClean(res)
}

func deepDiff(a, b interface{}, path string) []string {

	out := []string{}
	if a == nil && b == nil {
		return out
	}

	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return []string{path + fmt.Sprintf(" => different types [%v] and [%v]", a, b)}
	}

	switch a.(type) {

	case int32:
		if a.(int32) != b.(int32) {
			return []string{path + " => different int32 values"}
		}
		return out

	case int64:
		if a.(int64) != b.(int64) {
			return []string{path + " => different int64 values"}
		}
		return out

	case int:
		if a.(int) != b.(int) {
			return []string{path + " => different int values"}
		}
		return out

	case bool:
		if a.(bool) != b.(bool) {
			return []string{path + " => different bool values"}
		}
		return out

	case float64:
		if a.(float64) != b.(float64) {
			return []string{path + " => different float64 values"}
		}
		return out

	case string:
		if a.(string) != b.(string) {
			return []string{path + " => different string values"}
		}
		return out

	case map[string]interface{}:

		am := a.(map[string]interface{})
		bm := b.(map[string]interface{})

		for k, av := range am {
			p := path + "/" + k
			bv, find := bm[k]
			if !find {
				out = append(out, p+" => find only in A")
				continue
			}

			outLeft := deepDiff(av, bv, p)
			out = append(out, outLeft...)
		}

		for k := range bm {
			if _, find := am[k]; !find {
				out = append(out, path+"/"+k+" => find only in B")
				continue
			}
		}

		return out

	case []interface{}:

		am := a.([]interface{})
		bm := b.([]interface{})

		if len(am) != len(bm) {
			return []string{path + " => different array length"}
		}

		for i := range am {
			p := path + "/" + strconv.Itoa(i)
			outLeft := deepDiff(am[i], bm[i], p)
			if len(outLeft) > 0 {
				out = append(out, outLeft...)
			}
		}
		return out
	}

	return []string{path + " => //unknown error"}
}

func sortAndClean(res []string) []string {

	sort.Strings(res)

	if len(res) < 2 {
		return res
	}

	if len(res) == 2 && res[0] == res[1] {
		return []string{res[0]}
	}

	var i, j int
	var found bool
	for i, j = 0, 1; i < len(res) && j < len(res); i, j = i+1, j+1 {

		for j < len(res) {
			if res[i] == res[j] {
				found = true
				if j+1 < len(res) {
					j++
					continue
				}
			}
			break
		}

		if i+1 != j {
			res[i+1] = res[j]
		}
	}

	if !found {
		return res
	}
	return res[0:i]
}
