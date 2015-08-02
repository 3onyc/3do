package parser

import (
	"bufio"
	"fmt"
	"github.com/3onyc/3do/model"
	"github.com/3onyc/3do/util"
	"io"
)

func ToModels(t string, l *List) *model.TodoList {
	tl := &model.TodoList{
		Title:  t,
		Groups: []*model.TodoGroup{},
	}

	for _, g := range l.Groups {
		tg := &model.TodoGroup{
			Title: g.Title,
			Items: []*model.TodoItem{},
		}

		for _, i := range g.Items {
			ti := &model.TodoItem{
				Title:       i.Title,
				Done:        i.Done,
				Description: i.Description(),
			}

			tg.Items = append(tg.Items, ti)
		}

		tl.Groups = append(tl.Groups, tg)
	}

	return tl
}

func Read(t string, r io.Reader) (*model.TodoList, error) {
	c := NewContext(t)
	s := bufio.NewScanner(r)

	for s.Scan() {
		if err := s.Err(); err != nil {
			return nil, err
		}

		ReadLine(c, Line(s.Text()))
	}

	return ToModels(t, c.List), nil
}

func ReadLine(c *Context, l Line) {
	if c.ShouldSkip(l) {
		return
	}

	if l.IsGroupTitle() {
		if c.Item != nil {
			c.Item.RemoveTrailingSpace()
		}

		c.Item = nil
		c.Group = NewGroup(l.AsGroupTitle())
		c.List.AddGroup(c.Group)
	} else if c.Group != nil && l.IsItemTitle() {
		if c.Item != nil {
			c.Item.RemoveTrailingSpace()
		}

		c.Item = NewItem(l.AsItemTitle(), l.IsDone())
		c.Group.AddItem(c.Item)
	} else if c.Item != nil {
		// Skip whitespace preceding item description
		if len(c.Item.Lines) == 0 && l.Empty() {
			return
		}

		c.Item.AddLine(l)
		fmt.Printf("LNIS: %s\n", util.ShowNewLines(c.Item.Description()))
	}
}
