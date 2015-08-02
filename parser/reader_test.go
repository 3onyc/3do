package parser

import (
	"bytes"
	"github.com/3onyc/3do/model"
	"github.com/google/gofuzz"
	"github.com/kr/pretty"
	"reflect"
	//"strings"
	"testing"
)

var readerInput1 = `### Group 1

* Item 1

   Some content
   _which can be markdown_

* Item 2

### Group 2

* Item 1`

var readerOutput1 = &model.TodoList{
	Title: "Foo",
	Groups: []*model.TodoGroup{
		&model.TodoGroup{
			Title: "Group 1",
			Items: []*model.TodoItem{
				&model.TodoItem{
					Title:       "Item 1",
					Done:        false,
					Description: "Some content\n_which can be markdown_",
				},
				&model.TodoItem{
					Title: "Item 2",
					Done:  false,
				},
			},
		},
		&model.TodoGroup{
			Title: "Group 2",
			Items: []*model.TodoItem{
				&model.TodoItem{
					Title: "Item 1",
					Done:  false,
				},
			},
		},
	},
}

func TestReader(t *testing.T) {
	l, err := Read("Foo", bytes.NewBufferString(readerInput1))
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(l, readerOutput1) {
		t.Error("Output doesn't equal expected content")
		for _, d := range pretty.Diff(l, readerOutput1) {
			t.Errorf("% #v", d)
		}
	}

	//if len(l.Groups) < 2 {
	//	t.Errorf("Expected 2 groups got %d", len(l.Groups))
	//}

	//g1 := l.Groups[0]
	//if g1.Title != "Group 1" {
	//	t.Errorf("First group title should be 'Group 1', got '%s'", g1.Title)
	//}
	//g2 := l.Groups[1]
	//if g2.Title != "Group 2" {
	//	t.Errorf("Second group title should be 'Group 2', got '%s'", g2.Title)
	//}

	//if len(g1.Items) < 2 {
	//	t.Errorf("Expected 2 items got %d", len(g1.Items))
	//}
	//g1i1 := g1.Items[0]
	//if g1i1.Title != "Item 1" {
	//	t.Errorf("Item title should be 'Item 1', got '%s'", g1i1.Title)
	//}

	//expected := "Some content\n_which can be markdown_\n"
	//if g1i1.Description != expected {
	//	t.Errorf(
	//		"Item description should be '%s', got '%s'",
	//		strings.Replace(expected, "\n", "\\n\n", -1),
	//		strings.Replace(g1i1.Description, "\n", "\\n\n", -1),
	//	)
	//}

	//g2i1desc := l.Groups[1].Items[0].Description
	//if g2i1desc != "" {
	//	t.Errorf("Group 2 Item 1 description should be empty, is '%s'\n", g2i1desc)
	//}
}

func TestReaderFuzzing(t *testing.T) {
	fuzzer := fuzz.New()
	for i := 0; i < 1000; i++ {
		var title, input string
		fuzzer.Fuzz(&title)
		fuzzer.Fuzz(&input)

		Read(title, bytes.NewBufferString(input))
	}
}
