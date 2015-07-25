package dbtype

import (
	"database/sql"
	"encoding/json"
)

type NullInt64 struct {
	sql.NullInt64
}

func NewNullInt64(val int64) NullInt64 {
	return NullInt64{
		sql.NullInt64{
			Valid: true,
			Int64: val,
		},
	}
}

func (n *NullInt64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return json.Marshal(nil)
	} else {
		return json.Marshal(n.Int64)
	}
}
