package gforms_test

import (
	"html/template"

	. "launchpad.net/gocheck"

	"github.com/vmihailenco/gforms"
)

type FormTest struct{}

var _ = Suite(&FormTest{})

//------------------------------------------------------------------------------

type TestForm struct {
	*gforms.BaseForm
	Name *gforms.StringField
	Age  *gforms.Int64Field
}

func NewTestForm() *TestForm {
	f := &TestForm{
		BaseForm: &gforms.BaseForm{},
		Name:     gforms.NewStringField(),
		Age:      gforms.NewInt64Field(),
	}
	gforms.InitForm(f)
	return f
}

//------------------------------------------------------------------------------

func (t *TestForm) TestFormUsage(c *C) {
	f := NewTestForm()

	valueGetter := func(f gforms.Field) interface{} {
		bf := f.ToBaseField()
		switch bf.Name {
		case "Name":
			return "foo"
		case "Age":
			return "23"
		}
		panic("unreachable")
	}

	c.Assert(gforms.IsValid(f, valueGetter), Equals, false)

	c.Assert(f.Name.Value(), Equals, "foo")
	c.Assert(
		f.Name.Render(),
		Equals,
		template.HTML(`<input type="text" id="Name" name="Name" value="foo" />`),
	)

	c.Assert(f.Age.Value(), Equals, 23)
	c.Assert(
		f.Age.Render(),
		Equals,
		template.HTML(`<input type="text" id="Age" name="Age" value="23" />`),
	)
}
