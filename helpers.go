package gforms

import (
	"bytes"
	"fmt"
	"html/template"
	"path"
	"reflect"
)

var (
	WidgetTemplate, CheckboxTemplate, RadioTemplate *template.Template
	emptyHTML                                       = template.HTML("")
)

func newTemplate(filepath string) (*template.Template, error) {
	t := template.New(path.Base(filepath))
	t = t.Funcs(template.FuncMap{
		"renderField": RenderField,
		"renderLabel": RenderLabel,
		"renderError": RenderError,
	})
	return t.ParseFiles(filepath)
}

func init() {
	var err error

	WidgetTemplate, err = newTemplate("templates/gforms/widget.html")
	if err != nil {
		panic(err)
	}

	CheckboxTemplate, err = newTemplate("templates/gforms/checkbox.html")
	if err != nil {
		panic(err)
	}

	RadioTemplate, err = newTemplate("templates/gforms/radio.html")
	if err != nil {
		panic(err)
	}
}

func Render(field Field, attrs ...string) (template.HTML, error) {
	if reflect.ValueOf(field).IsNil() {
		return emptyHTML, nil
	}

	bf := field.ToBaseField()

	data := struct {
		Field     Field
		BaseField *BaseField
		Attrs     []string
		Radios    []template.HTML
	}{
		Field:     field,
		BaseField: bf,
		Attrs:     attrs,
	}

	var t *template.Template
	switch widget := bf.Widget.(type) {
	case *CheckboxWidget:
		t = CheckboxTemplate
	case *RadioWidget:
		data.Radios = widget.Radios(attrs, field.(SingleValueField).StringValue())
		t = RadioTemplate
	default:
		t = WidgetTemplate
	}

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, data); err != nil {
		return emptyHTML, err
	}

	return template.HTML(buf.String()), nil
}

func RenderError(f Field) (template.HTML, error) {
	bf := f.ToBaseField()
	if bf.ValidationError == nil {
		return emptyHTML, nil
	}
	error := fmt.Sprintf(`<span class="help-inline">%v</span>`, bf.ValidationError)
	return template.HTML(error), nil
}

func RenderLabel(f Field) (template.HTML, error) {
	bf := f.ToBaseField()
	label := fmt.Sprintf(`<label class="control-label" for="%v">%v</label>`, bf.Name, bf.Label)
	return template.HTML(label), nil
}

func RenderField(f Field, attrs []string) (template.HTML, error) {
	return f.Render(attrs...), nil
}
