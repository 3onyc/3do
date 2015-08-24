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
		if result := ShowNewLines(input); result != expected {
			t.Errorf("'%s' doesn't match expected '%s'", result, expected)
		}
	}
}

func TestTrimRightSpace(t *testing.T) {
	cases := map[string]string{
		" Foo \t\r\n": " Foo",
		"Bar   ":      "Bar",
	}

	for input, expected := range cases {
		if result := TrimRightSpace(input); result != expected {
			t.Errorf("'%s' doesn't match expected '%s'", result, expected)
		}
	}
}
