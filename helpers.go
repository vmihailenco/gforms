package gforms

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
)

var WidgetTemplate, CheckboxTemplate, RadioTemplate *template.Template
var emptyHTML = template.HTML("")

func newTemplate(path string) (*template.Template, error) {
	t := template.New(path)
	t = t.Funcs(template.FuncMap{
		"renderField": RenderField,
		"renderLabel": RenderLabel,
		"renderError": RenderError,
	})
	return t.ParseFiles(path)
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

func field(fIntrfc interface{}) (Field, error) {
	f, ok := fIntrfc.(Field)
	if !ok {
		return nil, errors.New("gforms: expected Field")
	}
	return f, nil
}

func Render(field Field, attrs ...string) (template.HTML, error) {
	context := map[string]interface{}{
		"field": field,
		"attrs": attrs,
	}

	bf := field.ToBaseField()

	var t *template.Template
	switch widget := bf.Widget.(type) {
	case *CheckboxWidget:
		t = CheckboxTemplate
	case *RadioWidget:
		context["radios"] = widget.Radios(attrs, field.(SingleValueField).StringValue())
		t = RadioTemplate
	default:
		t = WidgetTemplate
	}

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, context); err != nil {
		return emptyHTML, err
	}

	return template.HTML(buf.String()), nil
}

func RenderError(fI interface{}) (template.HTML, error) {
	f, err := field(fI)
	if err != nil {
		return emptyHTML, err
	}
	bf := f.ToBaseField()
	if bf.ValidationError == nil {
		return emptyHTML, nil
	}
	error := fmt.Sprintf(`<span class="help-inline">%v</span>`, bf.ValidationError)
	return template.HTML(error), nil
}

func RenderLabel(fI interface{}) (template.HTML, error) {
	f, err := field(fI)
	if err != nil {
		return emptyHTML, err
	}
	bf := f.ToBaseField()
	label := fmt.Sprintf(`<label class="control-label" for="%v">%v</label>`, bf.Name, bf.Label)
	return template.HTML(label), nil
}

func RenderField(fI interface{}, attrsI interface{}) (template.HTML, error) {
	f, err := field(fI)
	if err != nil {
		return emptyHTML, err
	}
	return f.Render(attrsI.([]string)...), nil
}
