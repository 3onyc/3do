package parser

import (
	"bufio"
	"fmt"
	"github.com/3onyc/3do/model"
	"github.com/3onyc/3do/util"
	"io"
	"strings"
)

type ParseError struct {
	Message string
}

func NewParseError(msg string) *ParseError {
	return &ParseError{
		Message: msg,
	}
}

func (p *ParseError) Error() string {
	return fmt.Sprintf("Parser Error: %s", p.Message)
}

type Token int

const (
	TOK_NONE Token = iota
	TOK_GROUP
	TOK_ITEM
	TOK_ITEM_DONE
	TOK_ITEM_DESC
)

func (t Token) String() string {
	switch t {
	case TOK_NONE:
		return "TOK_NONE"
	case TOK_GROUP:
		return "TOK_GROUP"
	case TOK_ITEM:
		return "TOK_ITEM"
	case TOK_ITEM_DONE:
		return "TOK_ITEM_DONE"
	case TOK_ITEM_DESC:
		return "TOK_ITEM_DESC"
	default:
		return "TOK_UNKNOWN"
	}
}

func GetToken(l string, prev Token) Token {
	switch {
	case len(l) > 4 && l[0:4] == "### ":
		return TOK_GROUP
	case len(l) > 2 && l[0:2] == "* ":
		if strings.HasPrefix(l[2:], "~~") && strings.HasSuffix(l, "~~") {
			return TOK_ITEM_DONE
		} else {
			return TOK_ITEM
		}
	case len(l) > 3 && l[0:3] == "   ":
		return TOK_ITEM_DESC
	default:
		return TOK_NONE
	}
}

func GetNormalised(l string, t Token) string {
	l = util.TrimRightSpace(l)

	switch t {
	case TOK_GROUP:
		return l[4:]
	case TOK_ITEM:
		return l[2:]
	case TOK_ITEM_DONE:
		return strings.TrimSuffix(strings.TrimPrefix(l, "* ~~"), "~~")
	case TOK_ITEM_DESC:
		if l == "" {
			return ""
		}
		return l[3:]
	default:
		return ""
	}
}

type Parser struct {
	Debug bool
}

func NewParser() *Parser {
	return &Parser{
		Debug: false,
	}
}

func (p *Parser) Parse(t string, r io.Reader) (*model.TodoList, error) {
	var (
		item  *model.TodoItem  = nil
		group *model.TodoGroup = nil
		list  *model.TodoList  = &model.TodoList{Title: t}

		s       *bufio.Scanner = bufio.NewScanner(r)
		prevTok Token          = TOK_NONE
	)

	for s.Scan() {
		l := s.Text()
		tok := GetToken(l, prevTok)
		norm := GetNormalised(l, tok)

		if p.Debug {
			fmt.Printf("PARSER | %14s | % #v\n", tok, norm)
		}

		switch tok {
		case TOK_GROUP:
			item = nil
			group = &model.TodoGroup{
				Title: norm,
			}
			list.Groups = append(list.Groups, group)
		case TOK_ITEM_DONE:
			fallthrough
		case TOK_ITEM:
			if group == nil {
				return nil, NewParseError(fmt.Sprintf("Encountered %s before TOK_GROUP", tok))
			}

			item = &model.TodoItem{
				Title: norm,
				Done:  tok == TOK_ITEM_DONE,
			}
			group.Items = append(group.Items, item)
		case TOK_ITEM_DESC:
			if item == nil {
				return nil, NewParseError("Encountered TOK_ITEM_DESC before TOK_ITEM")
			}

			item.Description = fmt.Sprintf("%s%s\n", item.Description, norm)
		}

		prevTok = tok
	}

	for _, g := range list.Groups {
		for _, i := range g.Items {
			i.Description = strings.TrimLeft(i.Description, "\n")
			i.Description = util.TrimRightSpace(i.Description)
		}
	}

	if s.Err() != nil {
		return nil, s.Err()
	}

	return list, nil
}

func Parse(t string, r io.Reader) (*model.TodoList, error) {
	return NewParser().Parse(t, r)
}
