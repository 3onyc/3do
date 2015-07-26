package model

import (
	"github.com/3onyc/3do/dbtype"
	"time"
)

type TodoGroup struct {
	ID        dbtype.NullInt64 `json:"id,string"`
	Title     string           `json:"title"`
	CreatedAt *time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time       `json:"updatedAt" db:"updated_at"`
	List      int64            `json:"list,string" db:"list_id"`
	Items     []int64          `json:"items,string"`
}
