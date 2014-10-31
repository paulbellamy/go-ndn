package tlv

import (
	"bytes"
	"encoding/binary"
	"io"
)

type UintTLV struct {
	T uint64
	V uint64
}

func (t UintTLV) Type() uint64 {
	return t.T
}

func (t UintTLV) WriteTo(w io.Writer) (int64, error) {
	buf := &bytes.Buffer{}
	err := WriteUint(buf, t.V)
	if err != nil {
		return 0, err
	}

	return ByteTLV{T: t.T, V: buf.Bytes()}.WriteTo(w)
}

func WriteUint(w io.Writer, v uint64) error {
	if v <= 0xff {
		return binary.Write(w, binary.BigEndian, uint8(v))
	} else if v <= 0xffff {
		return binary.Write(w, binary.BigEndian, uint16(v))
	} else if v <= 0xffffffff {
		return binary.Write(w, binary.BigEndian, uint32(v))
	} else {
		return binary.Write(w, binary.BigEndian, v)
	}
}

func (t UintTLV) MarshalBinary() ([]byte, error) {
	buf := &bytes.Buffer{}
	_, err := t.WriteTo(buf)
	return buf.Bytes(), err
}
