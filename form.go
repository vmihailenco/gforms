package gforms

import (
	"net/url"
	"reflect"
)

type FormValuer interface {
	FormValue(key string) string
}

type Form interface {
	SetErrors(map[string]error)
	Errors() map[string]error
}

func InitForm(form Form) {
	formStruct := reflect.ValueOf(form).Elem()
	typeOfForm := formStruct.Type()
	for i := 0; i < formStruct.NumField(); i++ {
		field, ok := formStruct.Field(i).Interface().(Field)
		if !ok {
			continue
		}

		name := typeOfForm.Field(i).Name
		field.SetName(name)
		field.SetLabel(name)
		field.Widget().Attrs().Set("id", name)
		field.Widget().Attrs().Set("name", name)
	}
}

func IsValid(f Form, data map[string][]interface{}) bool {
	s := reflect.ValueOf(f).Elem()
	errs := make(map[string]error, 0)
	for i := 0; i < s.NumField(); i++ {
		field, ok := s.Field(i).Interface().(Field)
		if !ok {
			continue
		}

		var value interface{}
		if _, ok := data[field.Name()]; ok {
			if field.IsMulti() {
				value = data[field.Name()]
			} else {
				value = data[field.Name()][0]
			}
		}

		if !IsFieldValid(field, value) {
			errs[field.Name()] = field.ValidationError()
		}
	}
	f.SetErrors(errs)

	return len(f.Errors()) == 0
}

func IsFormValid(f Form, data url.Values) bool {
	m := make(map[string][]interface{})
	for key, values := range data {
		valuesI := make([]interface{}, 0)
		for _, value := range values {
			valuesI = append(valuesI, value)
		}
		m[key] = valuesI
	}
	return IsValid(f, m)
}

type BaseForm struct {
	errors map[string]error
}

func (bf *BaseForm) SetErrors(errors map[string]error) {
	bf.errors = errors
}

func (bf *BaseForm) Errors() map[string]error {
	return bf.errors
}
