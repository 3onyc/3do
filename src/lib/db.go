package lib

import (
	"github.com/jmoiron/sqlx"
	"lib/model"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DB *sqlx.DB
)

func GetDB() *sqlx.DB {
	if DB == nil {
		DB = sqlx.MustConnect("sqlite3", "/tmp/3do.sqlite3")
		return DB
	}

	return DB
}

func CreateDBSchema(db *sqlx.DB) {
	db.MustExec(`
		CREATE TABLE IF NOT EXISTS "todo_items" (
			"id" integer,
			"title" varchar(255),
			"description" varchar(255),
			"done" bool,
			"done_at" datetime,
			"created_at" datetime,
			"updated_at" datetime,
			"group_id" integer ,

			 PRIMARY KEY ("id")
		 );
	`)
	db.MustExec(`
		CREATE INDEX IF NOT EXISTS idx_todo_items_group_id ON "todo_items"("group_id");
	`)

	db.MustExec(`
		CREATE TABLE IF NOT EXISTS "todo_groups" (
			"id" integer,
			"title" varchar(255),
			"created_at" datetime,
			"updated_at" datetime,
			"list_id" integer ,

			 PRIMARY KEY ("id")
		 );
	`)
	db.MustExec(`
		CREATE INDEX IF NOT EXISTS idx_todo_groups_list_id ON "todo_groups"("list_id");
	`)

	db.MustExec(`
		CREATE TABLE IF NOT EXISTS "todo_lists" (
			"id" integer,
			"title" varchar(255),
			"description" varchar(255),
			"created_at" datetime,
			"updated_at" datetime ,

			 PRIMARY KEY ("id")
		 );
	`)
}

func SeedDB(db *sqlx.DB) error {
	l1 := &model.TodoList{
		Title:       "Hello",
		Description: "Foo",
	}
	l2 := &model.TodoList{
		Title:       "Bye",
		Description: "Bar",
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

	g1 := &model.TodoGroup{
		Title: "Group 1",
	}

	if err := model.InsertTodoList(db, l1); err != nil {
		return err
	}
	if err := model.InsertTodoList(db, l2); err != nil {
		return err
	}

	if err := model.InsertTodoGroup(db, g1); err != nil {
		return err
	}

	i1.Group = g1.ID.Int64
	i2.Group = g1.ID.Int64
	if err := model.InsertTodoItem(db, i1); err != nil {
		return err
	}
	if err := model.InsertTodoItem(db, i2); err != nil {
		return err
	}

	return nil
}
