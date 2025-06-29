package models

import (
	"database/sql"
	"encoding/json"
)

type NullInt16 struct {
	sql.NullInt16
}

func (n *NullInt16) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Int16)
}

func (n *NullInt16) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		return nil
	}
	n.Valid = true
	return json.Unmarshal(data, &n.Int16)
}
