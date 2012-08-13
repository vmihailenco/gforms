package gforms

import (
	"fmt"
	"html/template"
	"strings"
	tTemplate "text/template"
)

//------------------------------------------------------------------------------

type Widget interface {
	Attrs() *WidgetAttrs
	Render([]string, ...string) template.HTML
}

type ChoiceWidget interface {
	Widget
	SetChoices(choices [][2]string)
}

//------------------------------------------------------------------------------

type BaseWidget struct {
	HTML  string
	attrs *WidgetAttrs
}

func (w *BaseWidget) Attrs() *WidgetAttrs {
	return w.attrs
}

func (w *BaseWidget) Render(attrs []string, values ...string) template.HTML {
	w.Attrs().FromSlice(attrs)
	html := fmt.Sprintf(w.HTML, w.Attrs().String(), tTemplate.HTMLEscapeString(values[0]))
	return template.HTML(html)
}

//------------------------------------------------------------------------------

type TextWidget struct {
	*BaseWidget
}

func NewTextWidget() *TextWidget {
	return &TextWidget{
		&BaseWidget{
			HTML: `<input%v value="%v" />`,
			attrs: &WidgetAttrs{
				attrs: [][2]string{{"type", "text"}},
			},
		},
	}
}

//------------------------------------------------------------------------------

type TextareaWidget struct {
	*BaseWidget
}

func NewTextareaWidget() *TextareaWidget {
	return &TextareaWidget{
		&BaseWidget{
			HTML: `<textarea%v>%v</textarea>`,
			attrs: &WidgetAttrs{
				attrs: make([][2]string, 0),
			},
		},
	}
}

//------------------------------------------------------------------------------

type CheckboxWidget struct {
	*BaseWidget
}

func NewCheckboxWidget() *CheckboxWidget {
	return &CheckboxWidget{
		&BaseWidget{
			HTML: `<input%v value="%v" />`,
			attrs: &WidgetAttrs{
				attrs: [][2]string{{"type", "checkbox"}},
			},
		},
	}
}

//------------------------------------------------------------------------------

type SelectWidget struct {
	*BaseWidget
	choices [][2]string
}

func (w *SelectWidget) SetChoices(choices [][2]string) {
	w.choices = choices
}

func (w *SelectWidget) Options(selValues ...string) []string {
	options := make([]string, 0, len(w.choices))
	for _, choice := range w.choices {
		value := tTemplate.HTMLEscapeString(choice[0])
		label := tTemplate.HTMLEscapeString(choice[1])
		attrs := ""
		for _, selValue := range selValues {
			if value == selValue {
				attrs = ` selected="selected"`
			}
		}
		option := fmt.Sprintf(`<option value="%v"%v>%v</option>`, value, attrs, label)
		options = append(options, option)
	}
	return options
}

func (w *SelectWidget) Render(attrs []string, values ...string) template.HTML {
	w.Attrs().FromSlice(attrs)
	options := strings.Join(w.Options(values...), "\n")
	selectHTML := fmt.Sprintf(w.HTML, w.Attrs().String(), options)
	return template.HTML(selectHTML)
}

func NewSelectWidget() *SelectWidget {
	return &SelectWidget{
		BaseWidget: &BaseWidget{
			HTML: `<select%v>%v</select>`,
			attrs: &WidgetAttrs{
				attrs: make([][2]string, 0),
			},
		},
	}
}

func NewMultiSelectWidget() *SelectWidget {
	return &SelectWidget{
		BaseWidget: &BaseWidget{
			HTML: `<select multiple="multiple"%v>%v</select>`,
			attrs: &WidgetAttrs{
				attrs: make([][2]string, 0),
			},
		},
	}
}

//------------------------------------------------------------------------------

type RadioWidget struct {
	*BaseWidget
	choices [][2]string
}

func (w *RadioWidget) SetChoices(choices [][2]string) {
	w.choices = choices
}

func (w *RadioWidget) Radios(attrs []string, checkedValue string) []template.HTML {
	id, _ := w.Attrs().Get("id")
	radios := make([]template.HTML, 0, len(w.choices))
	for i, choice := range w.choices {
		wAttrs := w.Attrs().Clone()
		wAttrs.Set("id", fmt.Sprintf("%v_%v", id, i))
		wAttrs.FromSlice(attrs)

		value := tTemplate.HTMLEscapeString(choice[0])
		label := tTemplate.HTMLEscapeString(choice[1])

		checked := ""
		if value == checkedValue {
			checked = ` checked="checked"`
		}

		radio := fmt.Sprintf(
			`<input%v value="%v"%v /> %v`,
			wAttrs.String(),
			value,
			checked,
			label)
		radios = append(radios, template.HTML(radio))
	}
	return radios
}

func (w *RadioWidget) Render(attrs []string, checkedValues ...string) template.HTML {
	panic("not implemented.")
	return template.HTML("")
}

func NewRadioWidget() *RadioWidget {
	return &RadioWidget{
		BaseWidget: &BaseWidget{
			attrs: &WidgetAttrs{
				attrs: [][2]string{{"type", "radio"}},
			},
		},
	}
}

//------------------------------------------------------------------------------

type FileWidget struct {
	*BaseWidget
}

func NewFileWidget() *FileWidget {
	return &FileWidget{
		&BaseWidget{
			HTML: `<input%v />`,
			attrs: &WidgetAttrs{
				attrs: [][2]string{{"type", "file"}},
			},
		},
	}
}

func (w *FileWidget) Render(attrs []string, values ...string) template.HTML {
	w.Attrs().FromSlice(attrs)
	html := fmt.Sprintf(w.HTML, w.Attrs().String())
	return template.HTML(html)
}
