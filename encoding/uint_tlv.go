package encoding

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

func (t UintTLV) WriteTo(w io.Writer) (n int64, err error) {
	buf := &bytes.Buffer{}

	if t.V <= 0xff {
		binary.Write(buf, binary.BigEndian, uint8(t.V))
	} else if t.V <= 0xffff {
		binary.Write(buf, binary.BigEndian, uint16(t.V))
	} else if t.V <= 0xffffffff {
		binary.Write(buf, binary.BigEndian, uint32(t.V))
	} else {
		binary.Write(buf, binary.BigEndian, t.V)
	}

	return ByteTLV{T: t.T, V: buf.Bytes()}.WriteTo(w)
}

func (t UintTLV) MarshalBinary() ([]byte, error) {
	buf := &bytes.Buffer{}
	_, err := t.WriteTo(buf)
	return buf.Bytes(), err
}
