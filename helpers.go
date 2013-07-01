package gforms

import (
	"bytes"
	"html/template"
	"os"
	"path"
	"reflect"
	"sync"
)

type templatesMap struct {
	L sync.RWMutex
	M map[string]*template.Template
}

var emptyHTML = template.HTML("")

var templates = &templatesMap{
	M: make(map[string]*template.Template),
}

var (
	WidgetTemplatePath   = rootDir() + "templates/gforms/widget.html"
	CheckboxTemplatePath = rootDir() + "templates/gforms/checkbox.html"
	RadioTemplatePath    = rootDir() + "templates/gforms/radio.html"
)

func rootDir() string {
	return os.Getenv("ROOT_DIR")
}

func getTemplate(filepath string) *template.Template {
	templates.L.RLock()
	t, ok := templates.M[filepath]
	if ok {
		templates.L.RUnlock()
		return t
	}
	templates.L.RUnlock()

	t = template.New(path.Base(filepath))
	t = t.Funcs(template.FuncMap{
		"field":       RenderField,
		"label":       RenderLabel,
		"field_error": RenderError,
	})
	var err error
	t, err = t.ParseFiles(filepath)
	if err != nil {
		panic(err)
	}

	templates.L.Lock()
	templates.M[filepath] = t
	templates.L.Unlock()

	return t
}

func RenderErrors(form Form) (template.HTML, error) {
	errors := form.Errors()
	if len(errors) == 0 {
		return emptyHTML, nil
	}

	fields := form.Fields()
	s := ""
	for name, e := range errors {
		field, ok := fields[name]
		if name == "" || (ok && field.Widget().IsHidden()) {
			s += `<div class="alert alert-error">` + e.Error() + `</div>` + "\n"
		}
	}
	return template.HTML(s), nil
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
	case *HiddenWidget:
		return RenderField(field, attrs)
	case *CheckboxWidget:
		t = getTemplate(CheckboxTemplatePath)
	case *RadioWidget:
		data.Radios = widget.Radios(attrs, field.(SingleValueField).StringValue())
		t = getTemplate(RadioTemplatePath)
	default:
		t = getTemplate(WidgetTemplatePath)
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

func RenderHiddenFields(form Form) (template.HTML, error) {
	fields := form.Fields()

	var html template.HTML
	for _, field := range fields {
		if field.Widget().IsHidden() {
			html += field.Render()
		}
	}
	return html, nil
}
