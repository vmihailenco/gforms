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
	Register((*TextareaStringField)(nil), func() interface{} {
		return NewTextareaStringField()
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
	SetFields(map[string]Field)
	Fields() map[string]Field

	SetErrors(map[string]error)
	Errors() map[string]error
}

func InitForm(form Form) error {
	formv := reflect.ValueOf(form).Elem()
	formt := formv.Type()
	tinfo := tinfoMap.TypeInfo(formt)

	fields := make(map[string]Field, len(tinfo.fields))
	for _, finfo := range tinfo.fields {
		fv := formv.FieldByIndex(finfo.idx)
		isNil := fv.IsNil()
		if isNil {
			fv.Set(reflect.ValueOf(finfo.constr()))
		}
		f := fv.Interface().(Field)
		if !f.HasName() {
			f.SetName(finfo.name)
		}
		if !f.HasLabel() {
			f.SetLabel(finfo.label)
		}
		if isNil {
			f.SetIsRequired(finfo.flags&fReq != 0)
		}
		fields[f.Name()] = f
	}
	form.SetFields(fields)

	return nil
}

type valueGetterFunc func(Field) interface{}

func IsValid(f Form, getValue valueGetterFunc) bool {
	formv := reflect.ValueOf(f).Elem()
	formt := formv.Type()
	tinfo := tinfoMap.TypeInfo(formt)

	errs := make(map[string]error, 0)
	for _, finfo := range tinfo.fields {
		fv := formv.FieldByIndex(finfo.idx)
		if fv.IsNil() {
			continue
		}

		f := fv.Interface().(Field)
		if !IsFieldValid(f, getValue(f)) {
			errs[f.Name()] = f.ValidationError()
		}
	}
	f.SetErrors(errs)

	return len(f.Errors()) == 0
}

func IsFormValid(form Form, formValues url.Values) bool {
	getValue := func(f Field) interface{} {
		if f.IsMultipart() {
			panic("IsFormValid() is called on multipart form (use IsMultipartFormValid())")
		} else {
			if f.IsMulti() {
				return formValues[f.Name()]
			} else {
				if values, ok := formValues[f.Name()]; ok {
					return values[0]
				}
			}
		}
		return nil
	}
	return IsValid(form, getValue)
}

func IsMultipartFormValid(form Form, multipartForm *multipart.Form) bool {
	getValue := func(f Field) interface{} {
		if f.IsMultipart() {
			if f.IsMulti() {
				return multipartForm.File[f.Name()]
			} else {
				if _, ok := multipartForm.File[f.Name()]; ok {
					return multipartForm.File[f.Name()][0]
				}
			}
		} else {
			if f.IsMulti() {
				return multipartForm.Value[f.Name()]
			} else {
				if values, ok := multipartForm.Value[f.Name()]; ok {
					return values[0]
				}
			}
		}
		return nil
	}
	return IsValid(form, getValue)
}

//------------------------------------------------------------------------------

type BaseForm struct {
	fields map[string]Field
	errors map[string]error
}

func (f *BaseForm) SetFields(fields map[string]Field) {
	f.fields = fields
}

func (f *BaseForm) Fields() map[string]Field {
	return f.fields
}

func (f *BaseForm) SetErrors(errors map[string]error) {
	f.errors = errors
}

func (f *BaseForm) Errors() map[string]error {
	return f.errors
}
