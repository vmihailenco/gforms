package gforms

import (
	"mime/multipart"
	"net/url"
	"reflect"
)

//------------------------------------------------------------------------------

func init() {
	Register((*StringField)(nil), func() interface{} {
		return NewStringField()
	})
	Register((*StringChoiceField)(nil), func() interface{} {
		return NewSelectStringField()
	})
	Register((*Int64Field)(nil), func() interface{} {
		return NewInt64Field()
	})
	Register((*Int64ChoiceField)(nil), func() interface{} {
		return NewSelectInt64Field()
	})
	Register((*BoolField)(nil), func() interface{} {
		return NewBoolField()
	})
	Register((*MultiStringChoiceField)(nil), func() interface{} {
		return NewMultiSelectStringField()
	})
	Register((*MultiInt64ChoiceField)(nil), func() interface{} {
		return NewMultiSelectInt64Field()
	})
}

//------------------------------------------------------------------------------

type Form interface {
	SetErrors(map[string]error)
	Errors() map[string]error
}

func InitForm(form Form) error {
	formv := reflect.ValueOf(form).Elem()
	formt := formv.Type()

	tinfo := tinfoMap.TypeInfo(formt)
	for _, finfo := range tinfo.fields {
		fv := formv.FieldByIndex(finfo.idx)
		if fv.IsNil() {
			fv.Set(reflect.ValueOf(finfo.constr()))
		}
		f := fv.Interface().(Field)

		bf := f.ToBaseField()
		if bf.Name == "" {
			bf.Name = finfo.name
		}
		if bf.Label == "" {
			bf.Label = finfo.name
		}

		attrs := bf.Widget.Attrs()
		if _, ok := attrs.Get("id"); !ok {
			attrs.Set("id", finfo.name)
		}
		if _, ok := attrs.Get("name"); !ok {
			attrs.Set("name", finfo.name)
		}

		bf.IsRequired = finfo.flags&fReq != 0
	}

	return nil
}

type valueGetterFunc func(Field) interface{}

func IsValid(f Form, getValue valueGetterFunc) bool {
	s := reflect.ValueOf(f).Elem()

	errs := make(map[string]error, 0)
	for i := 0; i < s.NumField(); i++ {
		field, ok := s.Field(i).Interface().(Field)
		if !ok {
			continue
		}
		bf := field.ToBaseField()

		if !IsFieldValid(field, getValue(field)) {
			errs[bf.Name] = bf.ValidationError
		}
	}
	f.SetErrors(errs)

	return len(f.Errors()) == 0
}

func IsFormValid(f Form, formValues url.Values) bool {
	getValue := func(field Field) (value interface{}) {
		bf := field.ToBaseField()

		if bf.IsMultipart {
			panic("IsFormValid() is called on multipart form (use IsMultipartFormValid())")
		} else {
			if bf.IsMulti {
				value = formValues[bf.Name]
			} else {
				if values, ok := formValues[bf.Name]; ok {
					value = values[0]
				}
			}
		}
		return
	}
	return IsValid(f, getValue)
}

func IsMultipartFormValid(f Form, multipartForm *multipart.Form) bool {
	getValue := func(field Field) (value interface{}) {
		bf := field.ToBaseField()

		if bf.IsMultipart {
			if bf.IsMulti {
				value = multipartForm.File[bf.Name]
			} else {
				if _, ok := multipartForm.File[bf.Name]; ok {
					value = multipartForm.File[bf.Name][0]
				}
			}
		} else {
			if bf.IsMulti {
				value = multipartForm.Value[bf.Name]
			} else {
				if values, ok := multipartForm.Value[bf.Name]; ok {
					value = values[0]
				}
			}
		}
		return
	}
	return IsValid(f, getValue)
}

//------------------------------------------------------------------------------

type BaseForm struct {
	errors map[string]error
}

func (bf *BaseForm) SetErrors(errors map[string]error) {
	bf.errors = errors
}

func (bf *BaseForm) Errors() map[string]error {
	return bf.errors
}
