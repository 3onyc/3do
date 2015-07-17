package model

import (
	"time"
)

type TodoItem struct {
	ID          uint      `json:"id,string"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Done        bool      `json:"done"`
	DoneAt      time.Time `json:"doneAt"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	GroupID     uint      `json:"group,string" sql:"index"`
}
