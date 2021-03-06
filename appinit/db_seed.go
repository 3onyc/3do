package appinit

import (
	"github.com/3onyc/3do/model"
	"github.com/3onyc/3do/util"
)

func SeedDB(ctx *util.Context) error {
	l1 := &model.TodoList{
		Title:       "Hello",
		Description: "Foo",
	}
	l2 := &model.TodoList{
		Title:       "Bye",
		Description: "Bar",
	}

	g1 := &model.TodoGroup{
		Title: "Group 1",
	}

	i1 := &model.TodoItem{
		Title:       "Item 1",
		Description: "# Foo",
		Done:        false,
	}
	i2 := &model.TodoItem{
		Title:       "Item 2",
		Description: "Bar",
		Done:        false,
	}

	if err := ctx.Lists.Insert(l1); err != nil {
		return err
	}
	if err := ctx.Lists.Insert(l2); err != nil {
		return err
	}

	g1.ListID = l1.ID.Int64
	if err := ctx.Groups.Insert(g1); err != nil {
		return err
	}

	i1.GroupID = g1.ID.Int64
	i2.GroupID = g1.ID.Int64
	if err := ctx.Items.Insert(i1); err != nil {
		return err
	}
	if err := ctx.Items.Insert(i2); err != nil {
		return err
	}

	return nil
}
