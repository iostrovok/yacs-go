package diff

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	. "gopkg.in/check.v1"
)

func loadFile(c *C, filename string) interface{} {

	file, err := os.Open(filename)
	c.Assert(err, IsNil)

	body, err := ioutil.ReadAll(file)
	c.Assert(err, IsNil)

	var s interface{}
	err = json.Unmarshal(body, &s)
	c.Assert(err, IsNil)

	return s
}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type diffTestSuite struct{}

var _ = Suite(&diffTestSuite{})

func (s *diffTestSuite) Test_SortUniq_V10(c *C) {
	// c.Skip("Skip")

	a := []string{}
	res := sortAndClean(a)
	c.Assert(res, DeepEquals, []string{})

	a = []string{"a"}
	res = sortAndClean(a)
	c.Assert(res, DeepEquals, []string{"a"})

}

func (s *diffTestSuite) Test_SortUniq_V90(c *C) {
	// c.Skip("Skip")

	a := []string{"a", "b"}
	res := sortAndClean(a)
	c.Assert(res, DeepEquals, []string{"a", "b"})

	a = []string{"b", "a"}
	res = sortAndClean(a)
	c.Assert(res, DeepEquals, []string{"a", "b"})

	a = []string{"b", "c", "a"}
	res = sortAndClean(a)
	c.Assert(res, DeepEquals, []string{"a", "b", "c"})

}

func (s *diffTestSuite) Test_SortUniq_V100(c *C) {
	// c.Skip("Skip")

	a := []string{}
	res := sortAndClean(a)
	c.Assert(res, DeepEquals, []string{})

	a = []string{"a", "a"}
	res = sortAndClean(a)
	c.Assert(res, DeepEquals, []string{"a"})

	a = []string{"b", "c", "a", "c", "c", "c", "c"}
	res = sortAndClean(a)
	c.Assert(res, DeepEquals, []string{"a", "b", "c"})

}

func (s *diffTestSuite) Test_SortUniq_V110(c *C) {
	// c.Skip("Skip")

	a := []string{}
	res := sortAndClean(a)
	c.Assert(res, DeepEquals, []string{})

}

func (s *diffTestSuite) Test_SortUniq_V120(c *C) {
	// c.Skip("Skip")

	a := []string{"a", "b", "c", "d"}
	res := sortAndClean(a)
	c.Assert(res, DeepEquals, []string{"a", "b", "c", "d"})

}

func (s *diffTestSuite) Test_SortUniq_V130(c *C) {
	// c.Skip("Skip")
	a := []string{"a", "b", "c", "d"}
	res := sortAndClean(a)
	c.Assert(res, DeepEquals, []string{"a", "b", "c", "d"})

	a = []string{"a", "b", "c", "d", "d"}
	res = sortAndClean(a)
	c.Assert(res, DeepEquals, []string{"a", "b", "c", "d"})

	a = []string{"a", "b", "c", "d", "d", "d"}
	res = sortAndClean(a)
	c.Assert(res, DeepEquals, []string{"a", "b", "c", "d"})

}

func (s *diffTestSuite) Test_SortUniq_V140(c *C) {
	// c.Skip("Skip")

	a := []string{"d", "b", "c", "a", "c", "b", "b", "c", "c", "c", "d"}
	res := sortAndClean(a)
	c.Assert(res, DeepEquals, []string{"a", "b", "c", "d"})
}

func (s *diffTestSuite) Test_VNil(c *C) {
	// c.Skip("Skip")

	res := Diff(nil, nil)
	c.Assert(res, DeepEquals, []string{})
}

func (s *diffTestSuite) Test_V1(c *C) {
	// c.Skip("Skip")
	a := loadFile(c, "./my.json")
	b := loadFile(c, "./my.json")
	res := Diff(a, b)

	// c.Assert(find, Equals, true)
	// c.Assert(len(res), Equals, 0)
	c.Assert(res, DeepEquals, []string{})
}

func (s *diffTestSuite) Test_V3(c *C) {
	// c.Skip("Skip")

	a := map[string]interface{}{
		"map": map[string]interface{}{
			"1": "1",
		},
		"list":   []interface{}{"1"},
		"string": "32",
		"bool":   false,
		"float":  32.03,
		"int32":  int32(32),
		"int64":  int64(32),
		"int":    int(32),

		"map-bad": map[string]interface{}{
			"2": "1",
			"1": "1",
		},
		"list-bad":   []interface{}{"1"},
		"string-bad": "64",
		"bool-bad":   false,
		"float-bad":  64.03,
		"int32-bad":  int32(64),
		"int64-bad":  int64(64),
		"int-bad":    int(64),
	}
	b := map[string]interface{}{
		"map": map[string]interface{}{
			"1": "1",
		},
		"list":   []interface{}{"1"},
		"string": "32",
		"bool":   false,
		"float":  32.03,
		"int32":  int32(32),
		"int64":  int64(32),
		"int":    int(32),

		"map-bad": map[string]interface{}{
			"2": "1",
			"3": "1",
		},
		"list-bad":   []interface{}{"1", "2"},
		"string-bad": "32",
		"bool-bad":   true,
		"float-bad":  32.03,
		"int32-bad":  int32(32),
		"int64-bad":  int64(32),
		"int-bad":    int(32),
	}

	res := Diff(a, b)

	resCheck := []string{
		"/bool-bad => different bool values",
		"/float-bad => different float64 values",
		"/int-bad => different int values",
		"/int32-bad => different int32 values",
		"/int64-bad => different int64 values",
		"/list-bad => different array length",
		"/map-bad/1 => find only in A",
		"/map-bad/3 => find only in B",
		"/string-bad => different string values",
	}

	c.Assert(res, DeepEquals, resCheck)
}

func (s *diffTestSuite) Test_V4(c *C) {
	// c.Skip("Skip")

	a := map[string]interface{}{
		"list": []interface{}{"1", 2, 2.1},
	}
	b := map[string]interface{}{
		"list": []interface{}{"1", "2", "2.1"},
	}
	res := Diff(a, b)

	resCheck := []string{"/list/1 => different types [2] and [2]", "/list/2 => different types [2.1] and [2.1]"}
	c.Assert(res, DeepEquals, resCheck)
}
