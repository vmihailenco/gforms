package gforms

import (
	"testing"
)

func TestAttrs(t *testing.T) {
	attrs := &WidgetAttrs{}

	attrs.Set("foo", "bar")
	if v, flag := attrs.Get("foo"); v != "bar" || flag != true {
		t.Errorf("Got (%v, %v), expected (bar, true).", v, flag)
	}

	attrs.Set("foo", "bar2")
	if v, flag := attrs.Get("foo"); v != "bar2" || flag != true {
		t.Errorf("Got (%v, %v), expected (bar2, true).", v, flag)
	}

	if v, flag := attrs.Pop("foo"); v != "bar2" || flag != true {
		t.Errorf("Got (%v, %v), expected (bar2, true).", v, flag)
	}

	if v, flag := attrs.Pop("foo"); v != "" || flag != false {
		t.Errorf(`Got (%v, %v), expected ("", false).`, v, flag)
	}

	if v, flag := attrs.Get("foo"); v != "" || flag != false {
		t.Errorf(`Got (%v, %v), expected ("", false).`, v, flag)
	}
}
