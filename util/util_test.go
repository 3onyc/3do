package util

import (
	"testing"
)

func TestShowNewLines(t *testing.T) {
	cases := map[string]string{
		"Foo":      "Foo",
		"Bar\nBaz": "Bar\\n\nBaz",
	}

	for input, expected := range cases {
		if ShowNewLines(input) != expected {
			t.Errorf("'%s' doesn't match expected '%s'")
		}
	}
}

func TestTrimRightSpace(t *testing.T) {
	cases := map[string]string{
		" Foo \t\r\n": " Foo",
		"Bar   ":      "Bar",
	}

	for input, expected := range cases {
		if TrimRightSpace(input) != expected {
			t.Errorf("'%s' doesn't match expected '%s'")
		}
	}
}
