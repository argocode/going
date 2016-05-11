package nulls

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"encoding/xml"
	"strconv"
)

// NullInt32 adds an implementation for int32
// that supports proper JSON encoding/decoding.
type NullInt32 struct {
	Int32 int32
	Valid bool // Valid is true if Int32 is not NULL
}

// NewNullInt32 returns a new, properly instantiated
// NullInt object.
func NewNullInt32(i int32) NullInt32 {
	return NullInt32{Int32: i, Valid: true}
}

// NewNullInt32Ptr returns a pointer to a new, properly instantiated
// NullInt object.
func NewNullInt32Ptr(i int32) *NullInt32 {
	return &NullInt32{Int32: i, Valid: true}
}

// Scan implements the Scanner interface.
func (ns *NullInt32) Scan(value interface{}) error {
	n := sql.NullInt64{Int64: int64(ns.Int32)}
	err := n.Scan(value)
	ns.Int32, ns.Valid = int32(n.Int64), n.Valid
	return err
}

// Value implements the driver Valuer interface.
func (ns NullInt32) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return int64(ns.Int32), nil
}

// MarshalJSON marshals the underlying value to a
// proper JSON representation.
func (ns NullInt32) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.Int32)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON will unmarshal a JSON value into
// the propert representation of that value.
func (ns *NullInt32) UnmarshalJSON(text []byte) error {
	txt := string(text)
	ns.Valid = true
	if txt == "null" {
		ns.Valid = false
		return nil
	}
	i, err := strconv.ParseInt(txt, 10, 32)
	if err != nil {
		ns.Valid = false
		return err
	}
	j := int32(i)
	ns.Int32 = j
	return nil
}

// UnmarshalXML will unmarshal an XML value into
// the proper representation of that value
func (ns *NullInt32) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	ns.Valid = true
	for _, attr := range start.Attr {
		if attr.Name.Local == "nil" {
			ns.Valid = false
			break
		}
	}
	return d.DecodeElement(&ns.Int32, &start)
}
