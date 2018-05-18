package loader

import (
	"testing"

	. "gopkg.in/check.v1"
)

var testFileJSON map[string]interface{}
var testFileContent []byte

func init() {
	testFileJSON = map[string]interface{}{
		"languageCode": "en",
	}

	testFileContent = []byte(`{"languageCode": "en"}`)
}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type loaderTestSuite struct{}

var _ = Suite(&loaderTestSuite{})

func (s *loaderTestSuite) Test_loadFile_V01(c *C) {
	body, err := loadFile("./my.json")
	c.Assert(err, IsNil)
	c.Assert(body, DeepEquals, testFileContent)
}

func (s *loaderTestSuite) Test_loadFile_V02(c *C) {
	body, err := loadFile("file://my.json")
	c.Assert(err, IsNil)
	c.Assert(body, DeepEquals, testFileContent)
}

func (s *loaderTestSuite) Test_isURL_V01(c *C) {
	c.Assert(isURL("file://my.json"), Equals, false)
	c.Assert(isURL("my.json"), Equals, false)
	c.Assert(isURL("./my.json"), Equals, false)
	c.Assert(isURL("file:///User/Ivan/my.json"), Equals, false)
	c.Assert(isURL("./my.json"), Equals, false)
}

func (s *loaderTestSuite) Test_isURL_V02(c *C) {
	c.Assert(isURL("http://google.com/my.json"), Equals, true)
	c.Assert(isURL("https://google.com/my.json"), Equals, true)
}

func (s *loaderTestSuite) Test_GetURI_V01(c *C) {
	body, err := GetURI("my.json")
	c.Assert(err, IsNil)
	c.Assert(body, DeepEquals, testFileJSON)
}

func (s *loaderTestSuite) Test_GetURI_V02(c *C) {
	body, err := GetURI("file://my.json")
	c.Assert(err, IsNil)
	c.Assert(body, DeepEquals, testFileJSON)
}
