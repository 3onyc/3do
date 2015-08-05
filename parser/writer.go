package parser

import (
	"fmt"
	"github.com/3onyc/3do/model"
	"io"
	"strings"
)

type Writer struct{}

func NewWriter() *Writer {
	return &Writer{}
}

func (w *Writer) Write(l *model.TodoList, wr io.Writer) error {
	for gi, g := range l.Groups {
		if gi > 0 {
			fmt.Fprint(wr, "\n\n")
		}
		fmt.Fprintf(wr, "### %s", g.Title)

		for _, i := range g.Items {
			if !i.Done {
				fmt.Fprintf(wr, "\n\n* %s", i.Title)
			} else {
				fmt.Fprintf(wr, "\n\n* ~~%s~~", i.Title)
			}

			if i.Description != "" {
				fmt.Fprint(wr, "\n")

				// Indent item description
				for _, l := range strings.Split(i.Description, "\n") {
					fmt.Fprintf(wr, "\n   %s", l)
				}
			}
		}
	}

	return nil
}
