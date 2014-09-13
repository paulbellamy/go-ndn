package encoding

import (
	"encoding/binary"
	"errors"
	"io"
)

var ErrUnexpectexUintTLVLengthValue = errors.New("Unexpected UintTLV length value")

func (t *UintTLV) ReadFrom(r io.Reader) (n int64, err error) {
	tlvType, n, err := ReadNumber(r)
	if err != nil {
		return n, err
	}
	t.Type = tlvType

	length, n2, err := ReadNumber(r)
	n += n2
	if err != nil {
		return n, err
	}

	switch length {
	case 1:
		var value uint8
		err = binary.Read(r, binary.BigEndian, &value)
		t.Value = uint64(value)
	case 2:
		var value uint16
		err = binary.Read(r, binary.BigEndian, &value)
		t.Value = uint64(value)
	case 4:
		var value uint32
		err = binary.Read(r, binary.BigEndian, &value)
		t.Value = uint64(value)
	case 8:
		var value uint64
		err = binary.Read(r, binary.BigEndian, &value)
		t.Value = uint64(value)
	default:
		err = ErrUnexpectexUintTLVLengthValue
	}

	if err == nil {
		n += int64(length)
	}

	return
}
