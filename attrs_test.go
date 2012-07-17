package gforms_test

import (
	. "launchpad.net/gocheck"

	"github.com/vmihailenco/gforms"
)

type AttrsTest struct{}

var _ = Suite(&AttrsTest{})

func (t *AttrsTest) TestAttrs(c *C) {
	attrs := &gforms.WidgetAttrs{}

	attrs.Set("foo", "bar")
	v, exists := attrs.Get("foo")
	c.Check(exists, Equals, true)
	c.Check(v, Equals, "bar")

	attrs.Set("foo", "bar2")
	v, exists = attrs.Get("foo")
	c.Check(exists, Equals, true)
	c.Check(v, Equals, "bar2")

	v, exists = attrs.Pop("foo")
	c.Check(exists, Equals, true)
	c.Check(v, Equals, "bar2")

	v, exists = attrs.Pop("foo")
	c.Check(exists, Equals, false)
	c.Check(v, Equals, "")

	v, exists = attrs.Get("foo")
	c.Check(exists, Equals, false)
	c.Check(v, Equals, "")
}
