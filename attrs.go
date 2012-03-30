package gforms

import (
	"fmt"
	"strings"
	tTemplate "text/template"
)

type WidgetAttrs struct {
	attrs [][2]string
}

func (w *WidgetAttrs) Clone() *WidgetAttrs {
	return &WidgetAttrs{
		attrs: w.attrs[:],
	}
}

func (w *WidgetAttrs) Set(name, value string) {
	value = tTemplate.HTMLEscapeString(value)

	exists := false
	for i := range w.attrs {
		attr := &w.attrs[i]
		if attr[0] == name {
			exists = true
			attr[1] = value
		}
	}
	if !exists {
		w.attrs = append(w.attrs, [...]string{name, value})
	}
}

func (w *WidgetAttrs) Get(name string) (string, bool) {
	for _, attr := range w.attrs {
		if attr[0] == name {
			return attr[1], true
		}
	}
	return "", false
}

func (w *WidgetAttrs) Pop(name string) (string, bool) {
	for i, attr := range w.attrs {
		if attr[0] == name {
			w.attrs = append(w.attrs[:i], w.attrs[i+1:]...)
			return attr[1], true
		}
	}
	return "", false
}

func (w *WidgetAttrs) Names() []string {
	names := make([]string, 0, len(w.attrs))
	for _, attr := range w.attrs {
		names = append(names, attr[0])
	}
	return names
}

func (w *WidgetAttrs) String() string {
	attrsArr := make([]string, 0, len(w.attrs))
	for _, attr := range w.attrs {
		attrsArr = append(attrsArr, fmt.Sprintf(`%v="%v"`, attr[0], attr[1]))
	}
	if len(attrsArr) > 0 {
		return " " + strings.Join(attrsArr, " ")
	}
	return ""
}

func (w *WidgetAttrs) FromSlice(attrs []string) {
	for i := 0; i < len(attrs); i += 2 {
		w.Set(attrs[i], attrs[i+1])
	}
}
