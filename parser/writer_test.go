package parser

import (
	"bytes"
	"github.com/3onyc/3do/model"
	"github.com/3onyc/3do/util"
	"testing"
)

var writerOutput1 = `### Group 1

* Item 1

   Foo
   Bar
   # Baz

* ~~Item 2~~

### Group 2

* Item 1`

var writerInput1 = &model.TodoList{
	Title:       "Foo",
	Description: "List Description",
	Groups: []*model.TodoGroup{
		&model.TodoGroup{
			Title: "Group 1",
			Items: []*model.TodoItem{
				&model.TodoItem{
					Title:       "Item 1",
					Done:        false,
					Description: "Foo\nBar\n# Baz",
				},
				&model.TodoItem{
					Title: "Item 2",
					Done:  true,
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

func TestWriter(t *testing.T) {
	buf := bytes.NewBufferString("")
	Write(writerInput1, buf)

	if buf.String() != writerOutput1 {
		t.Error("=== Expected ===")
		t.Error(util.ShowNewLines(writerOutput1))
		t.Error("===== Got =====")
		t.Error(util.ShowNewLines(buf.String()))
	}
}
