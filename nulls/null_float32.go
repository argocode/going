package nulls

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"encoding/xml"
	"strconv"
)

// NullFloat32 adds an implementation for float32
// that supports proper JSON encoding/decoding.
type NullFloat32 struct {
	Float32 float32
	Valid   bool // Valid is true if Float32 is not NULL
}

// NewNullFloat32 returns a new, properly instantiated
// NullFloat32 object.
func NewNullFloat32(i float32) NullFloat32 {
	return NullFloat32{Float32: i, Valid: true}
}

// NewNullFloat32Ptr returns a pointer to a new, properly instantiated
// NullFloat32 object.
func NewNullFloat32Ptr(i float32) *NullFloat32 {
	return &NullFloat32{Float32: i, Valid: true}
}

// Scan implements the Scanner interface.
func (ns *NullFloat32) Scan(value interface{}) error {
	n := sql.NullFloat64{Float64: float64(ns.Float32)}
	err := n.Scan(value)
	ns.Float32, ns.Valid = float32(n.Float64), n.Valid
	return err
}

// Value implements the driver Valuer interface.
func (ns NullFloat32) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return float64(ns.Float32), nil
}

// MarshalJSON marshals the underlying value to a
// proper JSON representation.
func (ns NullFloat32) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.Float32)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON will unmarshal a JSON value into
// the propert representation of that value.
func (ns *NullFloat32) UnmarshalJSON(text []byte) error {
	txt := string(text)
	ns.Valid = true
	if txt == "null" {
		ns.Valid = false
		return nil
	}
	i, err := strconv.ParseFloat(txt, 32)
	if err != nil {
		ns.Valid = false
		return err
	}
	j := float32(i)
	ns.Float32 = j
	return nil
}

// UnmarshalXML will unmarshal an XML value into
// the proper representation of that value
func (ns *NullFloat32) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	ns.Valid = true
	for _, attr := range start.Attr {
		if attr.Name.Local == "nil" {
			ns.Valid = false
			break
		}
	}
	return d.DecodeElement(&ns.Float32, &start)
}
