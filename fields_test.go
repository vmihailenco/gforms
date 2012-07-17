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

	c.Check(gforms.IsFieldValid(f, nil), Equals, false)
	c.Check(f.ValidationError.Error(), Equals, "This field is required.")
	c.Check(f.Value(), Equals, "")
	c.Check(f.Render(), Equals, template.HTML(`<input type="text" value="" />`))
}

func (t *FieldsTest) TestRequiredStringFieldPassValidation(c *C) {
	f := gforms.NewStringField()

	c.Check(gforms.IsFieldValid(f, "foo"), Equals, true)
	c.Check(f.ValidationError, IsNil)
	c.Check(f.Value(), Equals, "foo")
	c.Check(f.Render(), Equals, template.HTML(`<input type="text" value="foo" />`))
}

func (t *FieldsTest) TestOptionalStringFieldPassValidation(c *C) {
	f := gforms.NewStringField()
	f.IsRequired = false

	c.Check(gforms.IsFieldValid(f, nil), Equals, true)
	c.Check(f.ValidationError, IsNil)
	c.Check(f.Value(), Equals, "")
	c.Check(f.Render(), Equals, template.HTML(`<input type="text" value="" />`))
}

func (t *FieldsTest) TestSetName(c *C) {
	f := gforms.NewStringField()
	f.Name = "foo"
	f.Label = "fooLabel"

	c.Check(f.Render(), Equals, template.HTML(`<input type="text" value="" />`))
}

//------------------------------------------------------------------------------

func (t *FieldsTest) TestSelectStringFieldValidation(c *C) {
	f := gforms.NewSelectStringField()
	f.SetChoices([]gforms.StringChoice{{"foo", "bar"}})

	c.Check(gforms.IsFieldValid(f, "x"), Equals, false)
	c.Check(f.ValidationError.Error(), Equals, "x is not valid choice.")
	c.Check(f.Value(), Equals, "")

	c.Check(gforms.IsFieldValid(f, "foo"), Equals, true)
	c.Check(f.ValidationError, IsNil)
	c.Check(f.Value(), Equals, "foo")
}

func (t *FieldsTest) TestSelectInt64FieldValidation(c *C) {
	f := gforms.NewSelectInt64Field()
	f.SetChoices([]gforms.Int64Choice{{1, "foo"}})

	c.Check(gforms.IsFieldValid(f, 0), Equals, false)
	c.Check(f.ValidationError.Error(), Equals, "0 is not valid choice.")
	c.Check(f.Value(), Equals, int64(0))

	c.Check(gforms.IsFieldValid(f, 1), Equals, true)
	c.Check(f.ValidationError, IsNil)
	c.Check(f.Value(), Equals, int64(1))
}

//------------------------------------------------------------------------------

func (t *FieldsTest) TestMultiSelectStringFieldValidation(c *C) {
	f := gforms.NewMultiSelectStringField()
	f.SetChoices([]gforms.StringChoice{{"foo", "bar"}, {"go", "Golang"}})

	c.Check(gforms.IsFieldValid(f, []interface{}{"x"}), Equals, false)
	c.Check(f.ValidationError.Error(), Equals, "x is not valid choice.")
	c.Check(f.Value(), IsNil)

	c.Check(gforms.IsFieldValid(f, []interface{}{"foo"}), Equals, true)
	c.Check(f.ValidationError, IsNil)
	c.Check(f.Value(), DeepEquals, []string{"foo"})

	c.Check(gforms.IsFieldValid(f, []interface{}{"foo", "go"}), Equals, true)
	c.Check(f.ValidationError, IsNil)
	c.Check(f.Value(), DeepEquals, []string{"foo", "go"})
}

func (t *FieldsTest) TestMultiSelectInt64FieldValidation(c *C) {
	f := gforms.NewMultiSelectInt64Field()
	f.SetChoices([]gforms.Int64Choice{{1, "bar"}, {2, "Golang"}})

	c.Check(gforms.IsFieldValid(f, []interface{}{0}), Equals, false)
	c.Check(f.ValidationError.Error(), Equals, "0 is not valid choice.")
	c.Check(f.Value(), IsNil)

	c.Check(gforms.IsFieldValid(f, []interface{}{1}), Equals, true)
	c.Check(f.ValidationError, IsNil)
	c.Check(f.Value(), DeepEquals, []int64{1})

	c.Check(gforms.IsFieldValid(f, []interface{}{1, 2}), Equals, true)
	c.Check(f.ValidationError, IsNil)
	c.Check(f.Value(), DeepEquals, []int64{1, 2})
}
