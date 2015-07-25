package util

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type Context struct {
	Router *mux.Router
	DB     *sqlx.DB
}

func NewContext(r *mux.Router, db *sqlx.DB) *Context {
	return &Context{
		Router: r,
		DB:     db,
	}
}
