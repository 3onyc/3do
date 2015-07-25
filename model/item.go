package model

import (
	"github.com/3onyc/threedo-backend/dbtype"
	"time"
)

type TodoItem struct {
	ID          dbtype.NullInt64 `json:"id,string"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Done        bool             `json:"done"`
	DoneAt      *time.Time       `json:"doneAt" db:"done_at"`
	CreatedAt   *time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt   *time.Time       `json:"updatedAt" db:"updated_at"`
	Group       int64            `json:"group,string" db:"group_id"`
}
