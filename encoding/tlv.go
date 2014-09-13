package encoding

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var ErrUnexpectexUintTLVLengthValue = errors.New("Unexpected UintTLV length value")

type TLV struct {
	Type  uint64
	Value []byte
}

// Try to parse the tlv's bytes as a number, returning an error if it is not
// valid.
func (t *TLV) Uint() (result uint64, err error) {
	r := bytes.NewReader(t.Value)
	switch len(t.Value) {
	case 1:
		var value uint8
		err = binary.Read(r, binary.BigEndian, &value)
		result = uint64(value)
	case 2:
		var value uint16
		err = binary.Read(r, binary.BigEndian, &value)
		result = uint64(value)
	case 4:
		var value uint32
		err = binary.Read(r, binary.BigEndian, &value)
		result = uint64(value)
	case 8:
		var value uint64
		err = binary.Read(r, binary.BigEndian, &value)
		result = uint64(value)
	default:
		err = ErrUnexpectexUintTLVLengthValue
	}
	return
}
