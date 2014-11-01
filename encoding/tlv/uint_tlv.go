package tlv

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/paulbellamy/go-ndn/encoding"
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
	_, err := WriteUint(buf, t.V)
	if err != nil {
		return 0, err
	}

	return ByteTLV{T: t.T, V: buf.Bytes()}.WriteTo(w)
}

func (t UintTLV) MarshalBinary() ([]byte, error) {
	buf := &bytes.Buffer{}
	_, err := t.WriteTo(buf)
	return buf.Bytes(), err
}

func WriteUint(w io.Writer, v uint64) (n int64, err error) {
	if v <= 0xff {
		n = 1
		err = binary.Write(w, binary.BigEndian, uint8(v))
	} else if v <= 0xffff {
		n = 2
		err = binary.Write(w, binary.BigEndian, uint16(v))
	} else if v <= 0xffffffff {
		n = 4
		err = binary.Write(w, binary.BigEndian, uint32(v))
	} else {
		n = 8
		err = binary.Write(w, binary.BigEndian, v)
	}
	if err != nil {
		n = 0
	}
	return
}

func ReadUint(r io.Reader) (uint64, int64, error) {
	buffer := &bytes.Buffer{}
	n, err := io.CopyN(buffer, r, 8)
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return 0, n, err
	}

	switch n {
	case 1:
		var value uint8
		err = binary.Read(buffer, binary.BigEndian, &value)
		return uint64(value), n, err
	case 2:
		var value uint16
		err = binary.Read(buffer, binary.BigEndian, &value)
		return uint64(value), n, err
	case 4:
		var value uint32
		err = binary.Read(buffer, binary.BigEndian, &value)
		return uint64(value), n, err
	case 8:
		var value uint64
		err = binary.Read(buffer, binary.BigEndian, &value)
		return uint64(value), n, err
	default:
		return 0, n, &encoding.InvalidUnmarshalError{Message: "malformed uint value"}
	}
}
