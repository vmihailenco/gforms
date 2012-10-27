// +build appengine

package gforms

import (
	"fmt"
	"html/template"
	"net/url"

	"appengine/blobstore"
)

//------------------------------------------------------------------------------

func init() {
	Register((*BlobField)(nil), func() interface{} {
		return NewBlobField()
	})
}

//------------------------------------------------------------------------------

func IsBlobstoreFormValid(form Form, blobs map[string][]*blobstore.BlobInfo, formValues url.Values) bool {
	getValue := func(f Field) (value interface{}) {
		if f.IsMultipart() {
			if f.IsMulti() {
				value = blobs[f.Name()]
			} else {
				if _, ok := blobs[f.Name()]; ok {
					value = blobs[f.Name()][0]
				}
			}
		} else {
			if f.IsMulti() {
				value = formValues[f.Name()]
			} else {
				if values, ok := formValues[f.Name()]; ok {
					value = values[0]
				}
			}
		}
		return
	}
	return IsValid(form, getValue)
}

//------------------------------------------------------------------------------

type BlobField struct {
	*BaseField
}

func (f *BlobField) Value() *blobstore.BlobInfo {
	if f.iValue == nil {
		return nil
	}
	return f.iValue.(*blobstore.BlobInfo)
}

func (f *BlobField) Validate(rawValue interface{}) error {
	value, ok := rawValue.(*blobstore.BlobInfo)
	if !ok {
		return fmt.Errorf("Type %T is not supported.", rawValue)
	}

	if err := f.ApplyValidators(value); err != nil {
		return err
	}

	f.iValue = value
	return nil
}

func (f *BlobField) SetInitial(initial *blobstore.BlobInfo) {
	f.iValue = initial
}

func (f *BlobField) Render(attrs ...string) template.HTML {
	return f.Widget().Render(attrs)
}

func NewBlobField() *BlobField {
	return &BlobField{
		BaseField: &BaseField{
			widget:      NewFileWidget(),
			isMultipart: true,
		},
	}
}
