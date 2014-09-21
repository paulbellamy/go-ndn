package encoding

import (
	"bytes"
	"io"
)

type Reader struct {
	r io.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		r: r,
	}
}

func (r *Reader) Read() (TLV, error) {
	t, _, err := ReadNumber(r.r)
	if err != nil {
		return nil, err
	}

	length, _, err := ReadNumber(r.r)
	if err != nil {
		return nil, err
	}

	value := make([]byte, length)
	n, err := r.r.Read(value)
	if err != nil {
		return nil, err
	}

	if uint64(n) < length {
		return nil, io.ErrUnexpectedEOF
	}

	if r.isParentTLV(t) {
		children := []TLV{}
		childReader := NewReader(bytes.NewReader(value[0:n]))
		var child TLV
		var err error
		for child, err = childReader.Read(); err == nil; child, err = childReader.Read() {
			children = append(children, child)
		}
		if err != nil && err != io.EOF {
			return nil, err
		}

		return ParentTLV{
			T: t,
			V: children,
		}, nil
	} else {
		return ByteTLV{T: t, V: value[0:n]}, nil
	}
}

func (r *Reader) isParentTLV(t uint64) bool {
	switch t {
	case InterestType, DataType, NameType:
		return true
	default:
		return false
	}
}
