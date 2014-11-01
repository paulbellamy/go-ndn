package tlv

import (
	"fmt"
	"io"
	"reflect"
)

// Value is a value that can be stored in a TLV.
// It is either nil or an instance of one of these types:
//
//   uint64
//   []byte
//   []TLV
//   string
type Value interface{}

type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e UnsupportedTypeError) Error() string {
	return "unsupported type %s" + e.Type.String()
}

// IsValue reports whether v is a valid Value parameter type.
func IsValue(v interface{}) bool {
	if v == nil {
		return true
	}
	switch v.(type) {
	case uint64, []byte, []TLV, string:
		return true
	}
	return false
}

func ReadValue(r io.Reader, t uint64) (interface{}, int64, error) {
	// Figure out which type of value, and read that.
	valueReader, ok := valueReaders[t]
	if !ok {
		return nil, 0, fmt.Errorf("unsupported tlv type %d", t)
	}
	value := valueReader()
	n, err := value.ReadFrom(r)
	return value, n, err
}
