package parser

import (
	"bytes"
	"fmt"
	"strings"
)

type List struct {
	Title  string
	Groups []*Group
}

func NewList(t string) *List {
	return &List{
		Title: t,
	}
}

func (l *List) AddGroup(g *Group) {
	l.Groups = append(l.Groups, g)
}

type Group struct {
	Title string
	Items []*Item
}

func NewGroup(t string) *Group {
	return &Group{
		Title: t,
	}
}

func (g *Group) AddItem(i *Item) {
	g.Items = append(g.Items, i)
}

type Item struct {
	Title string
	Done  bool
	Lines []Line
}

func NewItem(t string, d bool) *Item {
	return &Item{
		Title: t,
		Done:  d,
	}
}

func (i *Item) AddLine(l Line) {
	i.Lines = append(i.Lines, l)
}

func (i *Item) RemoveTrailingSpace() {
	for len(i.Lines) > 0 && i.Lines[len(i.Lines)-1].Empty() {
		i.Lines = i.Lines[:len(i.Lines)-1]
	}
}

func (i *Item) Description() string {
	buf := bytes.NewBuffer(nil)
	lineCnt := len(i.Lines)
	for li, l := range i.Lines {
		if _, err := buf.WriteString(l.AsItemLine()); err != nil {
			fmt.Printf("Item.Description() | ERR | %s\n", err.Error())
			return buf.String()
		}

		if li != lineCnt-1 {
			if _, err := buf.WriteString("\n"); err != nil {
				fmt.Printf("Item.Description() | ERR | %s\n", err.Error())
				return buf.String()
			}
		}
	}

	return buf.String()
}

type Context struct {
	List  *List
	Group *Group
	Item  *Item
}

func NewContext(t string) *Context {
	return &Context{
		List: NewList(t),
	}
}

func (c *Context) ShouldSkip(l Line) bool {
	return l.Empty() && (c.Group == nil || c.Item == nil)
}

type Line string

func (l Line) String() string {
	return string(l)
}

func (l Line) TrimSpace() Line {
	return Line(strings.TrimSpace(l.String()))
}

func (l Line) Empty() bool {
	return l.TrimSpace() == ""
}

func (l Line) IsGroupTitle() bool {
	return len(l) > 4 && l[0:4] == "### "
}

func (l Line) AsGroupTitle() string {
	if !l.IsGroupTitle() {
		return ""
	}
	return l[4:].TrimSpace().String()
}

func (l Line) IsItemTitle() bool {
	return len(l) > 2 && l[0:2] == "* "
}

func (l Line) AsItemTitle() string {
	if !l.IsItemTitle() {
		return ""
	}
	return l[2:].TrimSpace().String()
}

func (l Line) AsItemLine() string {
	if l.Empty() || len(l) < 3 {
		return ""
	}
	return strings.TrimSpace(string(l[3:]))
}

func (l Line) IsDone() bool {
	return len(l) > 4 && l[2:4] == "~~"
}
