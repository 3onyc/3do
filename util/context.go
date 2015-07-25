package util

import (
	"github.com/3onyc/threedo-backend/model"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type Context struct {
	Router *mux.Router
	DB     *sqlx.DB
	Lists  *model.TodoListRepository
	Groups *model.TodoGroupRepository
	Items  *model.TodoItemRepository
}

func NewContext(r *mux.Router, db *sqlx.DB) *Context {
	return &Context{
		Router: r,
		DB:     db,
		Lists:  model.NewTodoListRepository(db),
		Groups: model.NewTodoGroupRepository(db),
		Items:  model.NewTodoItemRepository(db),
	}
}
