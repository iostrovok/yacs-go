package helper

import (
	"testing"

	. "gopkg.in/check.v1"
)

var mainGetLockNamesRes map[string]bool

func init() {
	mainGetLockNamesRes = map[string]bool{
		"header": true,
		"footer": true,
	}
}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type parentTestSuite struct{}

var _ = Suite(&parentTestSuite{})

func (s *parentTestSuite) Test_getLockNames_V01(c *C) {

	check := map[string]interface{}{
		"@lock_names": "header",
	}

	res := getLockNames(check)
	c.Assert(res, DeepEquals, map[string]bool{"header": true})
}

func (s *parentTestSuite) Test_getLockNames_V02(c *C) {
	// getLockNames(doc interface{}) map[string]bool

	check := map[string]interface{}{
		"@lock_names": []string{
			"header",
			"footer",
		},
	}

	res := getLockNames(check)
	c.Assert(res, DeepEquals, mainGetLockNamesRes)
}

func (s *parentTestSuite) Test_getLockNames_V03(c *C) {
	// getLockNames(doc interface{}) map[string]bool

	check := map[string]interface{}{
		"@lock_names": map[string]interface{}{
			"header": "1",
			"footer": "1",
		},
	}

	res := getLockNames(check)
	c.Assert(res, DeepEquals, mainGetLockNamesRes)
}

func (s *parentTestSuite) Test_getLockNames_V04(c *C) {
	// getLockNames(doc interface{}) map[string]bool

	check := map[string]interface{}{
		"@lock_names": map[string]interface{}{
			"header": map[string]interface{}{},
			"footer": true,
		},
	}

	res := getLockNames(check)
	c.Assert(res, DeepEquals, mainGetLockNamesRes)
}

/*
func (s *parentTestSuite) Test_loadFile_V02(c *C) {
	body, err := loadFile("file://my.json")
	c.Assert(err, IsNil)
	c.Assert(body, DeepEquals, testFileContent)
}

func (s *parentTestSuite) Test_isURL_V01(c *C) {
	c.Assert(isURL("file://my.json"), Equals, false)
	c.Assert(isURL("my.json"), Equals, false)
	c.Assert(isURL("./my.json"), Equals, false)
	c.Assert(isURL("file:///User/Ivan/my.json"), Equals, false)
	c.Assert(isURL("./my.json"), Equals, false)
}

func (s *parentTestSuite) Test_isURL_V02(c *C) {
	c.Assert(isURL("http://google.com/my.json"), Equals, true)
	c.Assert(isURL("https://google.com/my.json"), Equals, true)
}

func (s *parentTestSuite) Test_GetURI_V01(c *C) {
	body, err := GetURI("my.json")
	c.Assert(err, IsNil)
	c.Assert(body, DeepEquals, testFileJSON)
}

func (s *parentTestSuite) Test_GetURI_V02(c *C) {
	body, err := GetURI("file://my.json")
	c.Assert(err, IsNil)
	c.Assert(body, DeepEquals, testFileJSON)
}
*/
