package gforms

import (
	"bytes"
	"html/template"
	"path"
	"reflect"
)

var (
	WidgetTemplate   = newTemplate("templates/gforms/widget.html")
	CheckboxTemplate = newTemplate("templates/gforms/checkbox.html")
	RadioTemplate    = newTemplate("templates/gforms/radio.html")

	emptyHTML = template.HTML("")
)

func newTemplate(filepath string) *template.Template {
	t := template.New(path.Base(filepath))
	t = t.Funcs(template.FuncMap{
		"renderField": RenderField,
		"renderLabel": RenderLabel,
		"renderError": RenderError,
	})
	var err error
	t, err = t.ParseFiles(filepath)
	if err != nil {
		panic(err)
	}
	return t
}

func Render(field Field, attrs ...string) (template.HTML, error) {
	if reflect.ValueOf(field).IsNil() {
		return emptyHTML, nil
	}

	data := struct {
		Field  Field
		Attrs  []string
		Radios []template.HTML
	}{
		Field: field,
		Attrs: attrs,
	}

	var t *template.Template
	switch widget := field.Widget().(type) {
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
	err := f.ValidationError()
	if err == nil {
		return emptyHTML, nil
	}
	s := `<span class="help-inline">` + err.Error() + `</span>`
	return template.HTML(s), nil
}

func RenderLabel(f Field) (template.HTML, error) {
	label := f.Label()
	if label == "" {
		return emptyHTML, nil
	}
	if f.IsRequired() {
		label += "*"
	}
	s := `<label class="control-label" for="` + f.Name() + `">` + label + `</label>`
	return template.HTML(s), nil
}

func RenderField(f Field, attrs []string) (template.HTML, error) {
	return f.Render(attrs...), nil
}
