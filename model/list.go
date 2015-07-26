package model

import (
	"github.com/3onyc/3do/dbtype"
	"time"
)

type TodoList struct {
	ID          dbtype.NullInt64 `json:"id,string"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	CreatedAt   *time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt   *time.Time       `json:"updatedAt" db:"updated_at"`
	Groups      []int64          `json:"groups,string"`
}
