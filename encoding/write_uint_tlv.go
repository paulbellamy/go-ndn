package encoding

import (
	"encoding/binary"
	"io"
)

func WriteUintTLV(w io.Writer, tlv *UintTLV) error {
	err := WriteNumber(w, tlv.Type)
	if err != nil {
		return err
	}

	var length uint8
	var value interface{}
	if tlv.Value <= 0xff {
		length = 1
		value = uint8(tlv.Value)
	} else if tlv.Value <= 0xffff {
		length = 2
		value = uint16(tlv.Value)
	} else if tlv.Value <= 0xffffffff {
		length = 4
		value = uint32(tlv.Value)
	} else {
		length = 8
		value = uint64(tlv.Value)
	}

	err = binary.Write(w, binary.BigEndian, length)
	if err != nil {
		return err
	}

	return binary.Write(w, binary.BigEndian, value)
}

func (t *UintTLV) WriteTo(w io.Writer) error {
	return WriteUintTLV(w, t)
}
