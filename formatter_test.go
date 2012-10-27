package gforms

import (
	. "launchpad.net/gocheck"
)

type FormatterTest struct{}

var _ = Suite(&FormatterTest{})

func (t *FormatterTest) TestSplitWords(c *C) {
	table := []struct {
		s     string
		words []string
	}{
		{"", []string{}},
		{"FooBar", []string{"Foo", "Bar"}},
		{"HTTP", []string{"HTTP"}},
		{"HTTPReq", []string{"HTTP", "Req"}},
		{"HTTPReqX", []string{"HTTP", "Req", "X"}},
	}

	for _, row := range table {
		c.Assert(splitWords(row.s), DeepEquals, row.words)
	}
}
