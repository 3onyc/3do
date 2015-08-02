package parser

import (
	"fmt"
	"github.com/3onyc/3do/model"
	"io"
	"strings"
)

type Writer struct{}

func (wr *Writer) WriteItem(i *model.TodoItem, w io.Writer) {
	// Item Title
	if i.Done {
		fmt.Fprintf(w, "* ~~%s~~", i.Title)
	} else {
		fmt.Fprintf(w, "* %s", i.Title)
	}

	if strings.TrimSpace(i.Description) == "" {
		return
	}

	fmt.Fprint(w, "\n\n")
	lines := strings.Split(i.Description, "\n")
	lineCnt := len(lines)
	for i, l := range lines {
		if strings.TrimSpace(l) != "" {
			fmt.Fprintf(w, "   %s", l)
		}

		if lineCnt-1 != i {
			fmt.Fprint(w, "\n")
		}
	}
}

func (wr *Writer) WriteList(l *model.TodoList, w io.Writer) {
	groupCnt := len(l.Groups)

	for gi, g := range l.Groups {
		itemCnt := len(g.Items)
		fmt.Fprintf(w, "### %s", g.Title)

		if itemCnt == 0 && gi == groupCnt-1 {
			continue
		}

		fmt.Fprint(w, "\n\n")
		for ii, i := range g.Items {
			wr.WriteItem(i, w)
			if ii != itemCnt-1 {
				fmt.Fprint(w, "\n\n")
			}
		}

		if gi != groupCnt-1 {
			fmt.Fprint(w, "\n\n")
		}
	}
}

func Write(l *model.TodoList, w io.Writer) {
	(&Writer{}).WriteList(l, w)
}
