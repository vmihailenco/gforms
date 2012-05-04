package gforms

import (
	"html/template"
	"testing"
)

func TestRequiredStringFieldDoNotPassValidation(t *testing.T) {
	f := NewStringField()

	if IsFieldValid(f, nil) {
		t.Errorf("Field passed validation.")
	}

	expectedErr := "This field is required."
	if err := f.ValidationError; err.Error() != expectedErr {
		t.Errorf("Expected %v, got %v.", expectedErr, err)
	}

	expectedV := ""
	if v := f.Value(); v != expectedV {
		t.Errorf("Expected %v, got %v.", expectedV, v)
	}

	expectedHTML := template.HTML(`<input type="text" value="" />`)
	if html := f.Render(); html != expectedHTML {
		t.Errorf("Expected %v, got %v.", expectedHTML, html)
	}
}

func TestRequiredStringFieldPassValidation(t *testing.T) {
	f := NewStringField()

	if !IsFieldValid(f, "foo") {
		t.Errorf("Field did not pass validation: %v.", f.ValidationError)
	}

	if v := f.Value(); v != "foo" {
		t.Errorf("Expected foo, got %v", v)
	}

	expectedHTML := template.HTML(`<input type="text" value="foo" />`)
	if html := f.Render(); html != expectedHTML {
		t.Errorf("Expected %v, got %v.", expectedHTML, html)
	}
}

func TestOptionalStringFieldPassValidation(t *testing.T) {
	f := NewStringField()
	f.IsRequired = false

	if !IsFieldValid(f, nil) {
		t.Errorf("Field did not pass validation: %v.", f.ValidationError)
	}

	expectedHTML := template.HTML(`<input type="text" value="" />`)
	if html := f.Render(); html != expectedHTML {
		t.Errorf("Expected %v, got %v.", expectedHTML, html)
	}
}

func TestSetName(t *testing.T) {
	f := NewStringField()

	f.Name = "foo"
	if name := f.Name; name != "foo" {
		t.Errorf(`Expected "foo", got %v`, name)
	}

	f.Label = "fooLabel"
	if label := f.Label; label != "fooLabel" {
		t.Errorf(`Expected "fooLabel", got %v`, label)
	}

	expectedHTML := template.HTML(`<input type="text" value="" />`)
	if html := f.Render(); html != expectedHTML {
		t.Errorf("Expected %v, got %v.", expectedHTML, html)
	}
}

func TestSelectStringFieldValidation(t *testing.T) {
	f := NewSelectStringField()
	f.SetChoices([]StringChoice{{"foo", "bar"}})

	if IsFieldValid(f, "x") {
		t.Errorf("Field passed validation.")
	}

	expectedErr := "x is not valid choice."
	if err := f.ValidationError; err.Error() != expectedErr {
		t.Errorf(`Expected %v, got %v.`, expectedErr, err)
	}

	if !IsFieldValid(f, "foo") {
		t.Errorf(`Field did not pass validation: %v.`, f.ValidationError)
	}

	if v := f.Value(); v != "foo" {
		t.Errorf(`Expected "foo", got %v.`, v)
	}
}

func TestSelectInt64FieldValidation(t *testing.T) {
	f := NewSelectInt64Field()
	f.SetChoices([]Int64Choice{{1, "foo"}})

	if IsFieldValid(f, 0) {
		t.Errorf("Field passed validation.")
	}

	expectedErr := "0 is not valid choice."
	if err := f.ValidationError; err.Error() != expectedErr {
		t.Errorf(`Expected %v, got %v.`, expectedErr, err)
	}

	if !IsFieldValid(f, 1) {
		t.Errorf(`Field did not pass validation: %v.`, f.ValidationError)
	}

	if v := f.Value(); v != 1 {
		t.Errorf(`Expected 1, got %v.`, v)
	}
}

func TestMultiSelectStringFieldValidation(t *testing.T) {
	f := NewMultiSelectStringField()
	f.SetChoices([]StringChoice{{"foo", "bar"}, {"go", "Golang"}})

	if IsFieldValid(f, []interface{}{"x"}) {
		t.Errorf("Field passed validation.")
	}

	expectedErr := "x is not valid choice."
	if err := f.ValidationError; err.Error() != expectedErr {
		t.Errorf("Expected %v, got %v.", expectedErr, err)
	}

	if !IsFieldValid(f, []interface{}{"foo", "go"}) {
		t.Errorf("Field did not pass validation: %v.", f.ValidationError)
	}

	if v := f.Value(); len(v) != 2 || v[0] != "foo" || v[1] != "go" {
		t.Errorf(`Expected {"foo", "bar"}, got %v.`, v)
	}
}

func TestMultiSelectInt64FieldValidation(t *testing.T) {
	f := NewMultiSelectInt64Field()
	f.SetChoices([]Int64Choice{{1, "bar"}, {2, "Golang"}})

	if IsFieldValid(f, []interface{}{0}) {
		t.Errorf("Field passed validation.")
	}

	expectedErr := "0 is not valid choice."
	if err := f.ValidationError; err.Error() != expectedErr {
		t.Errorf("Expected %v, got %v.", expectedErr, err)
	}

	if !IsFieldValid(f, []interface{}{1, 2}) {
		t.Errorf("Field did not pass validation: %v.", f.ValidationError)
	}

	if v := f.Value(); len(v) != 2 || v[0] != 1 || v[1] != 2 {
		t.Errorf(`Expected {1, 2}, got %v.`, v)
	}
}
