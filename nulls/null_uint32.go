package nulls

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"encoding/xml"
	"strconv"
)

// NullUInt32 adds an implementation for int
// that supports proper JSON encoding/decoding.
type NullUInt32 struct {
	UInt32 uint32
	Valid  bool // Valid is true if Int is not NULL
}

// NewNullUInt32 returns a new, properly instantiated
// NullInt object.
func NewNullUInt32(i uint32) NullUInt32 {
	return NullUInt32{UInt32: i, Valid: true}
}

// NewNullUInt32Ptr returns a new, properly instantiated
// NullInt object.
func NewNullUInt32Ptr(i uint32) *NullUInt32 {
	return &NullUInt32{UInt32: i, Valid: true}
}

// Scan implements the Scanner interface.
func (ns *NullUInt32) Scan(value interface{}) error {
	n := sql.NullInt64{Int64: int64(ns.UInt32)}
	err := n.Scan(value)
	ns.UInt32, ns.Valid = uint32(n.Int64), n.Valid
	return err
}

// Value implements the driver Valuer interface.
func (ns NullUInt32) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return int64(ns.UInt32), nil
}

// MarshalJSON marshals the underlying value to a
// proper JSON representation.
func (ns NullUInt32) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.UInt32)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON will unmarshal a JSON value into
// the propert representation of that value.
func (ns *NullUInt32) UnmarshalJSON(text []byte) error {
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
	j := uint32(i)
	ns.UInt32 = j
	return nil
}

// UnmarshalXML will unmarshal an XML value into
// the proper representation of that value
func (ns *NullUInt32) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	ns.Valid = true
	for _, attr := range start.Attr {
		if attr.Name.Local == "nil" {
			ns.Valid = false
			break
		}
	}
	return d.DecodeElement(&ns.UInt32, &start)
}
