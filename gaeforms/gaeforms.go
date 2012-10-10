package gaeforms

import (
	"fmt"
	"html/template"
	"net/url"

	"appengine/blobstore"

	"github.com/vmihailenco/gforms"
)

//------------------------------------------------------------------------------

func init() {
	gforms.Register((*BlobField)(nil), func() interface{} {
		return NewBlobField()
	})
}

//------------------------------------------------------------------------------

func IsBlobstoreFormValid(f gforms.Form, blobs map[string][]*blobstore.BlobInfo, formValues url.Values) bool {
	getValue := func(field gforms.Field) (value interface{}) {
		bf := field.ToBaseField()
		if bf.IsMultipart {
			if bf.IsMulti {
				value = blobs[bf.Name]
			} else {
				if _, ok := blobs[bf.Name]; ok {
					value = blobs[bf.Name][0]
				}
			}
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
	return gforms.IsValid(f, getValue)
}

//------------------------------------------------------------------------------

type BlobField struct {
	*gforms.BaseField
}

func (f *BlobField) Value() *blobstore.BlobInfo {
	if f.IValue == nil {
		return nil
	}
	return f.IValue.(*blobstore.BlobInfo)
}

func (f *BlobField) Validate(rawValue interface{}) error {
	value, ok := rawValue.(*blobstore.BlobInfo)
	if !ok {
		return fmt.Errorf("Type %T is not supported.", rawValue)
	}

	if err := f.ApplyValidators(value); err != nil {
		return err
	}

	f.IValue = value
	return nil
}

func (f *BlobField) SetInitial(initial *blobstore.BlobInfo) {
	f.IValue = initial
}

func (f *BlobField) Render(attrs ...string) template.HTML {
	return f.Widget.Render(attrs)
}

func NewBlobField() *BlobField {
	return &BlobField{
		BaseField: &gforms.BaseField{
			Widget:      gforms.NewFileWidget(),
			IsMultipart: true,
			IsRequired:  true,
		},
	}
}
