package parser

import (
	"testing"
)

func TestLineIsEmpty(t *testing.T) {
	cases := map[string]bool{
		"\t":  true,
		"\n":  true,
		"\r":  true,
		"Foo": false,
	}

	for input, expect := range cases {
		if result := Line(input).Empty(); result != expect {
			t.Errorf("result should be '%b', got '%b'", expect, result)
		}
	}
}

func TestLineTrimSpace(t *testing.T) {
	cases := map[string]string{
		" Foo":    "Foo",
		"Bar ":    "Bar",
		"\rBaz\n": "Baz",
	}

	for input, expect := range cases {
		if result := Line(input).TrimSpace().String(); result != expect {
			t.Errorf("result should be '%s', got '%s'", expect, result)
		}
	}
}

func TestLineToString(t *testing.T) {
	expect := "foo"
	if result := Line("foo").String(); result != expect {
		t.Errorf("result should be '%s', is '%s' instead", expect, result)
	}
}

func TestLineIsGroupTitle(t *testing.T) {
	cases := map[string]bool{
		"### Foo":     true,
		"###NotTitle": false,
		"NotTitle":    false,
	}

	for input, expect := range cases {
		if result := Line(input).IsGroupTitle(); result != expect {
			t.Errorf("result should be '%b', got '%b'", expect, result)
		}
	}
}

func TestLineAsGroupTitle(t *testing.T) {
	cases := map[string]string{
		"### Foo":     "Foo",
		"###NotTitle": "",
		"NotTitle":    "",
	}

	for input, expect := range cases {
		if result := Line(input).AsGroupTitle(); result != expect {
			t.Errorf("result should be '%s', got '%s'", expect, result)
		}
	}
}

func TestLineIsItemTitle(t *testing.T) {
	cases := map[string]bool{
		"* Foo":     true,
		"*NotTitle": false,
		"NotTitle":  false,
	}

	for input, expect := range cases {
		if result := Line(input).IsItemTitle(); result != expect {
			t.Errorf("result should be '%b', got '%b'", expect, result)
		}
	}
}

func TestLineAsItemTitle(t *testing.T) {
	cases := map[string]string{
		"* Foo":     "Foo",
		"*NotTitle": "",
		"NotTitle":  "",
	}

	for input, expect := range cases {
		if result := Line(input).AsItemTitle(); result != expect {
			t.Errorf("result should be '%s', got '%s'", expect, result)
		}
	}
}

func TestAsItemLine(t *testing.T) {
	cases := map[string]string{
		"   # Line": "# Line",
		"   Line":   "Line",
		"   ":       "",
		" X":        "", // Out of slice bounds
	}

	for input, expect := range cases {
		if result := Line(input).AsItemLine(); result != expect {
			t.Errorf("result should be '%s', got '%s'", expect, result)
		}
	}
}

func TestItemLineDone(t *testing.T) {
	cases := map[string]bool{
		"* ~~Done~~": true,
		"* NotDone":  false,
		" X":         false, // Out of slice bounds
	}

	for input, expect := range cases {
		if result := Line(input).IsDone(); result != expect {
			t.Errorf("result should be '%b', got '%b'", expect, result)
		}
	}
}

func TestItemRemoveTrailingSpace(t *testing.T) {
	type ItemResult struct {
		Input    *Item
		Expected int
	}

	emptyDesc := NewItem("Empty Desc", false)
	emptyDesc.Lines = []Line{
		"",
	}

	lineDesc := NewItem("Line Desc", false)
	lineDesc.Lines = []Line{
		"   Line", "", "",
	}

	cases := []ItemResult{
		{emptyDesc, 0},
		{lineDesc, 1},
	}

	for _, c := range cases {
		i := c.Input
		i.RemoveTrailingSpace()
		if result := len(i.Lines); result != c.Expected {
			t.Errorf("result should be '%d', got '%d'", c.Expected, result)
		}
	}
}

func TestItemDescription(t *testing.T) {
	lineDesc := NewItem("Line Desc", false)
	lineDesc.Lines = []Line{
		"   Line",
	}
	twoLineDesc := NewItem("Line Desc", false)
	twoLineDesc.Lines = []Line{
		"   Line",
		"   Line 2",
	}

	cases := map[string]*Item{
		"Line":         lineDesc,
		"Line\nLine 2": twoLineDesc,
	}

	for expect, input := range cases {
		if result := input.Description(); result != expect {
			t.Errorf("result should be '%s', got '%s'", expect, result)
		}
	}
}
