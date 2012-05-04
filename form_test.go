package gforms

import (
	"html/template"
	"testing"
)

type TestForm struct {
	*BaseForm
	Name *StringField
	Age  *Int64Field
}

func TestFormUsage(t *testing.T) {
	f := &TestForm{
		BaseForm: &BaseForm{},
		Name:     NewStringField(),
		Age:      NewInt64Field(),
	}
	InitForm(f)

	valueGetter := func(f Field) interface{} {
		bf := f.ToBaseField()
		switch bf.Name {
		case "Name":
			return "foo"
		case "Age":
			return "23"
		}
		return "missing"
	}

	if !IsValid(f, valueGetter) {
		t.Errorf("Form did not pass validation: %v.", f.Errors())
	}

	if v := f.Name.Value(); v != "foo" {
		t.Errorf("Expected foo, got %v.", v)
	}
	expectedHTML := template.HTML(`<input type="text" id="Name" name="Name" value="foo" />`)
	if html := f.Name.Render(); html != expectedHTML {
		t.Errorf("Expected %v, got %v.", expectedHTML, html)
	}

	if v := f.Age.Value(); v != 23 {
		t.Errorf("Expected 23, got %v.", v)
	}
	expectedHTML = template.HTML(`<input type="text" id="Age" name="Age" value="23" />`)
	if html := f.Age.Render(); html != expectedHTML {
		t.Errorf("Expected %v, got %v.", expectedHTML, html)
	}
}
