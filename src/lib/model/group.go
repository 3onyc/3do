package model

import (
	"time"
)

type TodoGroup struct {
	ID        uint       `json:"id,string"`
	Title     string     `json:"title"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	ListID    uint       `json:"list,string" sql:"index"`
	Items     []TodoItem `json:"items,string"`
}
