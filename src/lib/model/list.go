package model

import (
	"time"
)

type TodoList struct {
	ID          uint        `json:"id,string"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	Groups      []TodoGroup `json:"groups"`
}
