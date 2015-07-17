package util

import (
	"database/sql"
	"encoding/json"
)

type NullInt64 struct {
	sql.NullInt64
}

func (n *NullInt64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return json.Marshal(nil)
	} else {
		return json.Marshal(n.Int64)
	}
}
