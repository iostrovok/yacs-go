package helper

import (
	"path/filepath"
	"strings"
)

// Context is just internal sturture.
type Context struct {
	uri string
}

func newContext() *Context {
	return &Context{}
}

func (c *Context) getDir(URI string) string {

	if filepath.IsAbs(URI) {
		return URI
	}

	if URI == "." {
		return c.uri
	}

	if strings.Index(URI, "#") == 0 {
		return filepath.Clean(filepath.Join(c.uri, URI))
	}

	dir, _ := filepath.Split(c.uri)
	return filepath.Clean(filepath.Join(dir, URI))
}

func (c *Context) setURI(URI string) {
	c.uri = URI
}

func (c *Context) copy() *Context {
	return &Context{
		uri: c.uri,
	}
}
