package parser

import (
	"bytes"
	"github.com/3onyc/3do/model"
	"github.com/google/gofuzz"
	"github.com/kr/pretty"
	"reflect"
	"testing"
)

var parserInput1 = `### Group 1

* Item 1

    Some content
   _which can be markdown_

* Item 2

* ~~Item 3~~

   Foo

### Group 2

* Item 1`

var parserOutput1 = &model.TodoList{
	Title: "Foo",
	Groups: []*model.TodoGroup{
		{
			Title: "Group 1",
			Items: []*model.TodoItem{
				{
					Title:       "Item 1",
					Done:        false,
					Description: " Some content\n_which can be markdown_",
				},
				{
					Title: "Item 2",
					Done:  false,
				},
				{
					Title:       "Item 3",
					Description: "Foo",
					Done:        true,
				},
			},
		},
		{
			Title: "Group 2",
			Items: []*model.TodoItem{
				{
					Title: "Item 1",
					Done:  false,
				},
			},
		},
	},
}

var parserInput2 = `### Group 1

* Item 1

   # Foo

* Item 2

   bar`

var parserOutput2 = &model.TodoList{
	Title: "Foo",
	Groups: []*model.TodoGroup{
		{
			Title: "Group 1",
			Items: []*model.TodoItem{
				{
					Title:       "Item 1",
					Description: "# Foo",
				},
				{
					Title:       "Item 2",
					Description: "bar",
				},
			},
		},
	},
}

var mangledInput1 = `* Invalid Item

   Invalid Desc

### Group 1

* Item 1

   Desc Line 1

### Group 2

   Desc Line 2`

var mangledInput2 = `### Group 1

* Item 1

   Desc

### Group 2

   Invalid Desc`

func TestParser(t *testing.T) {
	testParserCase(parserInput1, parserOutput1, nil, t)
	testParserCase(parserInput2, parserOutput2, nil, t)
	testParserCase(mangledInput1, nil, &ParseError{}, t)
	testParserCase(mangledInput2, nil, &ParseError{}, t)
}

func testParserCase(i string, o *model.TodoList, expectErr error, t *testing.T) {
	l, err := NewParser().Parse("Foo", bytes.NewBufferString(i))
	if err != nil && expectErr != nil {
		if _, ok := err.(*ParseError); !ok {
			t.Errorf("Expected ParseError, got % #v", err)
		}
	} else if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(o, l) {
		t.Errorf("% #v\n", pretty.Formatter(l))
		for _, d := range pretty.Diff(o, l) {
			t.Errorf("% #v\n", d)
		}
	}
}

func TestParserFuzzing(t *testing.T) {
	fuzzer := fuzz.New()
	for i := 0; i < 1000; i++ {
		var title, input string
		fuzzer.Fuzz(&title)
		fuzzer.Fuzz(&input)

		NewParser().Parse(title, bytes.NewBufferString(input))
	}
}
