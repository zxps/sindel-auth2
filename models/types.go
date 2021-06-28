package models

import (
	"database/sql"
	"encoding/json"
)

// NullString struct
type NullString struct {
	sql.NullString
}

// NullInt32 struct
type NullInt32 struct {
	sql.NullInt32
}

// Params struct
type Params struct {
	sql.NullString
}

// MarshalJSON for null string
func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid && len(ns.String) > 0 {
		return json.Marshal(ns.String)
	}

	return json.Marshal(nil)
}

// MarshalJSON for null string
func (v NullInt32) MarshalJSON() ([]byte, error) {
	if v.Valid && v.Int32 > 0 {
		return json.Marshal(v.Int32)
	}

	return json.Marshal(nil)
}

// MarshalJSON for params
func (p Params) MarshalJSON() ([]byte, error) {
	if p.Valid && len(p.String) > 0 {
		var result map[string]interface{}

		json.Unmarshal([]byte(p.String), &result)
		return json.Marshal(result)
	}

	return json.Marshal(map[string]interface{}{})
}

