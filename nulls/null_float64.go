package nulls

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"encoding/xml"
	"strconv"
)

// NullFloat64 replaces sql.NullFloat64 with an implementation
// that supports proper JSON encoding/decoding.
type NullFloat64 sql.NullFloat64

// NewNullFloat64 returns a new, properly instantiated
// NullFloat64 object.
func NewNullFloat64(i float64) NullFloat64 {
	return NullFloat64{Float64: i, Valid: true}
}

// NewNullFloat64Ptr returns a pointer to a new, properly instantiated
// NullFloat64 object.
func NewNullFloat64Ptr(i float64) *NullFloat64 {
	return &NullFloat64{Float64: i, Valid: true}
}

// Scan implements the Scanner interface.
func (ns *NullFloat64) Scan(value interface{}) error {
	n := sql.NullFloat64{Float64: ns.Float64}
	err := n.Scan(value)
	ns.Float64, ns.Valid = n.Float64, n.Valid
	return err
}

// Value implements the driver Valuer interface.
func (ns NullFloat64) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.Float64, nil
}

// MarshalJSON marshals the underlying value to a
// proper JSON representation.
func (ns NullFloat64) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.Float64)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON will unmarshal a JSON value into
// the propert representation of that value.
func (ns *NullFloat64) UnmarshalJSON(text []byte) error {
	t := string(text)
	ns.Valid = true
	if t == "null" {
		ns.Valid = false
		return nil
	}
	i, err := strconv.ParseFloat(t, 64)
	if err != nil {
		ns.Valid = false
		return err
	}
	ns.Float64 = i
	return nil
}

// UnmarshalXML will unmarshal an XML value into
// the proper representation of that value
func (ns *NullFloat64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	ns.Valid = true
	for _, attr := range start.Attr {
		if attr.Name.Local == "nil" {
			ns.Valid = false
			break
		}
	}
	return d.DecodeElement(&ns.Float64, &start)
}
