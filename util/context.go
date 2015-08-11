package util

import (
	"github.com/3onyc/3do/model"
	"github.com/asaskevich/EventBus"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type Context struct {
	Bus    *EventBus.EventBus
	DB     *sqlx.DB
	Logger log.Logger
	Router *mux.Router

	Lists  *model.TodoListRepository
	Groups *model.TodoGroupRepository
	Items  *model.TodoItemRepository
}

func NewContext(r *mux.Router, db *sqlx.DB, l log.Logger) *Context {
	return &Context{
		Bus:    EventBus.New(),
		DB:     db,
		Logger: l,
		Router: r,
		Lists:  model.NewTodoListRepository(db),
		Groups: model.NewTodoGroupRepository(db),
		Items:  model.NewTodoItemRepository(db),
	}
}

func (ctx *Context) Log(pairs ...interface{}) {
	ctx.Logger.Log(pairs...)
}
