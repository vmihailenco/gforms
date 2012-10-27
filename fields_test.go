package gforms_test

import (
	"html/template"

	. "launchpad.net/gocheck"

	"github.com/vmihailenco/gforms"
)

type FieldsTest struct{}

var _ = Suite(&FieldsTest{})

func (t *FieldsTest) TestRequiredStringFieldDoNotPassValidation(c *C) {
	f := gforms.NewStringField()
	f.SetIsRequired(true)

	c.Assert(gforms.IsFieldValid(f, nil), Equals, false)
	c.Assert(f.ValidationError().Error(), Equals, "This field is required")
	c.Assert(f.Value(), Equals, "")
	c.Assert(f.Render(), Equals, template.HTML(`<input type="text" value="" />`))
}

func (t *FieldsTest) TestRequiredStringFieldPassValidation(c *C) {
	f := gforms.NewStringField()

	c.Assert(gforms.IsFieldValid(f, "foo"), Equals, true)
	c.Assert(f.ValidationError(), IsNil)
	c.Assert(f.Value(), Equals, "foo")
	c.Assert(f.Render(), Equals, template.HTML(`<input type="text" value="foo" />`))
}

func (t *FieldsTest) TestOptionalStringFieldPassValidation(c *C) {
	f := gforms.NewStringField()

	c.Assert(gforms.IsFieldValid(f, nil), Equals, true)
	c.Assert(f.ValidationError(), IsNil)
	c.Assert(f.Value(), Equals, "")
	c.Assert(f.Render(), Equals, template.HTML(`<input type="text" value="" />`))
}

func (t *FieldsTest) TestSetName(c *C) {
	f := gforms.NewStringField()
	f.SetName("foo")
	f.SetLabel("fooLabel")

	c.Assert(
		f.Render(),
		Equals,
		template.HTML(`<input type="text" id="foo" name="foo" value="" />`),
	)
}

//------------------------------------------------------------------------------

func (t *FieldsTest) TestSelectStringFieldValidation(c *C) {
	f := gforms.NewSelectStringField()
	f.SetChoices([]gforms.StringChoice{{"foo", "bar"}})

	c.Assert(gforms.IsFieldValid(f, "x"), Equals, false)
	c.Assert(f.ValidationError().Error(), Equals, "x is invalid choice")
	c.Assert(f.Value(), Equals, "")

	c.Assert(gforms.IsFieldValid(f, "foo"), Equals, true)
	c.Assert(f.ValidationError(), IsNil)
	c.Assert(f.Value(), Equals, "foo")
}

func (t *FieldsTest) TestSelectInt64FieldValidation(c *C) {
	f := gforms.NewSelectInt64Field()
	f.SetChoices([]gforms.Int64Choice{{1, "foo"}})
	f.SetIsRequired(true)

	c.Assert(gforms.IsFieldValid(f, 0), Equals, false)
	c.Assert(f.ValidationError().Error(), Equals, "This field is required")
	c.Assert(f.Value(), Equals, int64(0))

	c.Assert(gforms.IsFieldValid(f, 1), Equals, true)
	c.Assert(f.ValidationError(), IsNil)
	c.Assert(f.Value(), Equals, int64(1))

	c.Assert(gforms.IsFieldValid(f, 2), Equals, false)
	c.Assert(f.ValidationError().Error(), Equals, "2 is invalid choice")
	c.Assert(f.Value(), Equals, int64(0))
}

//------------------------------------------------------------------------------

func (t *FieldsTest) TestMultiSelectStringFieldValidation(c *C) {
	f := gforms.NewMultiSelectStringField()
	f.SetChoices([]gforms.StringChoice{{"foo", "bar"}, {"go", "Golang"}})

	c.Assert(gforms.IsFieldValid(f, []interface{}{"x"}), Equals, false)
	c.Assert(f.ValidationError().Error(), Equals, "x is invalid choice")
	c.Assert(f.Value(), IsNil)

	c.Assert(gforms.IsFieldValid(f, []interface{}{"foo"}), Equals, true)
	c.Assert(f.ValidationError(), IsNil)
	c.Assert(f.Value(), DeepEquals, []string{"foo"})

	c.Assert(gforms.IsFieldValid(f, []interface{}{"foo", "go"}), Equals, true)
	c.Assert(f.ValidationError(), IsNil)
	c.Assert(f.Value(), DeepEquals, []string{"foo", "go"})
}

func (t *FieldsTest) TestMultiSelectInt64FieldValidation(c *C) {
	f := gforms.NewMultiSelectInt64Field()
	f.SetChoices([]gforms.Int64Choice{{1, "bar"}, {2, "Golang"}})

	c.Assert(gforms.IsFieldValid(f, []interface{}{0}), Equals, false)
	c.Assert(f.ValidationError().Error(), Equals, "0 is invalid choice")
	c.Assert(f.Value(), IsNil)

	c.Assert(gforms.IsFieldValid(f, []interface{}{1}), Equals, true)
	c.Assert(f.ValidationError(), IsNil)
	c.Assert(f.Value(), DeepEquals, []int64{1})

	c.Assert(gforms.IsFieldValid(f, []interface{}{1, 2}), Equals, true)
	c.Assert(f.ValidationError(), IsNil)
	c.Assert(f.Value(), DeepEquals, []int64{1, 2})
}
