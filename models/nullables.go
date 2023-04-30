package models

import (
	"database/sql"
	"encoding/json"
)

// Nullable Bool that overrides sql.NullBool
type NullBool struct {
	sql.NullBool
}

func (nb NullBool) MarshalJSON() ([]byte, error) {
	if nb.Valid {
		return json.Marshal(nb.Bool)
	}
	return json.Marshal(nil)
}

func (nb *NullBool) UnmarshalJSON(data []byte) error {
	var b *bool
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	if b != nil {
		nb.Valid = true
		nb.Bool = *b
	} else {
		nb.Valid = false
	}
	return nil
}

// Nullable Float64 that overrides sql.NullFloat64
type NullFloat64 struct {
	sql.NullFloat64
}

func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if nf.Valid {
		return json.Marshal(nf.Float64)
	}
	return json.Marshal(nil)
}

func (nf *NullFloat64) UnmarshalJSON(data []byte) error {
	var f *float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	if f != nil {
		nf.Valid = true
		nf.Float64 = *f
	} else {
		nf.Valid = false
	}
	return nil
}

// Nullable Int16 that overrides sql.NullInt16
type NullInt16 struct {
	sql.NullInt16
}

func (ni NullInt16) MarshalJSON() ([]byte, error) {
	if ni.Valid {
		return json.Marshal(ni.Int16)
	}
	return json.Marshal(nil)
}

func (ni *NullInt16) UnmarshalJSON(data []byte) error {
	var i *int16
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	if i != nil {
		ni.Valid = true
		ni.Int16 = *i
	} else {
		ni.Valid = false
	}
	return nil
}

// Nullable String that overrides sql.NullString
type NullString struct {
	sql.NullString
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		ns.Valid = true
		ns.String = *s
	} else {
		ns.Valid = false
	}
	return nil
}

type Optional[T any] struct {
	Defined bool
	Value   *T
}

// UnmarshalJSON is implemented by deferring to the wrapped type (T).
// It will be called only if the value is defined in the JSON payload.
func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	o.Defined = true
	return json.Unmarshal(data, &o.Value)
}
