package tlv

import (
	"bytes"
	"encoding"
	"errors"
	"io"
	"reflect"
)

var ErrUnexpectexUintTLVLengthValue = errors.New("Unexpected UintTLV length value")

type TLV interface {
	Type() uint64
	io.WriterTo
	encoding.BinaryMarshaler
}

type GenericTLV struct {
	T uint64
	V interface{}
}

func (t GenericTLV) Type() uint64 {
	return t.T
}

func (t GenericTLV) Bytes() []byte {
	return t.V.([]byte)
}

func (t GenericTLV) Children() []TLV {
	return t.V.([]TLV)
}

func (t GenericTLV) Uint() uint64 {
	return t.V.(uint64)
}

func (t GenericTLV) WriteTo(w io.Writer) (n int64, err error) {
	value := &bytes.Buffer{}
	_, err = writeValue(value, t.V)
	if err != nil {
		return
	}

	n, err = writeType(w, t.T)
	if err != nil {
		return
	}

	written, err := writeLength(w, value.Bytes())
	n += written
	if err != nil {
		return
	}

	written, err = writeValue(w, value.Bytes())
	n += written
	return
}

func writeType(w io.Writer, t uint64) (int64, error) {
	return WriteNumber(w, t)
}

func writeLength(w io.Writer, value []byte) (int64, error) {
	return WriteNumber(w, uint64(len(value)))
}

func writeValue(w io.Writer, value interface{}) (int64, error) {
	switch value := value.(type) {
	case uint64:
		return WriteUint(w, value)
	case []byte:
		written, err := w.Write(value)
		return int64(written), err
	case []TLV:
		var n int64
		for _, t := range value {
			written, err := t.WriteTo(w)
			n += written
			if err != nil {
				return n, err
			}
		}
		return n, nil
	case string:
		written, err := w.Write([]byte(value))
		return int64(written), err
	default:
		return 0, &UnsupportedTypeError{Type: reflect.TypeOf(value)}
	}
}

func (t GenericTLV) MarshalBinary() ([]byte, error) {
	buf := &bytes.Buffer{}
	_, err := t.WriteTo(buf)
	return buf.Bytes(), err
}

/*
// Try to parse the tlv's bytes as a number, returning an error if it is not
// valid.
func ReadUint(b []byte) (result uint64, err error) {
	r := bytes.NewReader(b)
	switch len(b) {
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
*/
