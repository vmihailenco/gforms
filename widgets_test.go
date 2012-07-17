package gforms_test

import (
	"html/template"

	. "launchpad.net/gocheck"

	"github.com/vmihailenco/gforms"
)

type WidgetsTest struct{}

var _ = Suite(&WidgetsTest{})

func (t *WidgetsTest) TestWidgets(c *C) {
	textWidget := gforms.NewTextWidget()
	checkboxWidget := gforms.NewCheckboxWidget()

	var widgetTests = []struct {
		given, expected template.HTML
	}{
		{
			textWidget.Render(nil, ""),
			template.HTML(`<input type="text" value="" />`),
		},
		{
			textWidget.Render([]string{"name", "foo"}, ""),
			template.HTML(`<input type="text" name="foo" value="" />`),
		},
		{
			checkboxWidget.Render(nil, ""),
			template.HTML(`<input type="checkbox" value="" />`),
		},
	}

	for _, tt := range widgetTests {
		c.Check(tt.given, Equals, tt.expected)
	}
}
