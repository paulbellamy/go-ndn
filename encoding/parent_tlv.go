package encoding

import (
	"bytes"
	"io"
)

func ParentTLV(t uint64, children ...TLV) TLV {
	return &parentTLV{
		T: t,
		V: children,
	}
}

type parentTLV struct {
	T uint64
	V []TLV
}

func (t *parentTLV) Type() uint64 {
	return t.T
}

func (t *parentTLV) WriteTo(w io.Writer) (n int64, err error) {
	childrenBytes, err := t.marshalChildrenToBinary()
	if err != nil {
		return
	}

	n, err = WriteNumber(w, t.T)
	if err != nil {
		return
	}

	written, err := WriteNumber(w, uint64(len(childrenBytes)))
	n += written
	if err != nil {
		return
	}

	written2, err := w.Write(childrenBytes)
	n += int64(written2)
	return
}

func (t *parentTLV) MarshalBinary() ([]byte, error) {
	buf := &bytes.Buffer{}
	_, err := t.WriteTo(buf)
	return buf.Bytes(), err
}

func (t *parentTLV) marshalChildrenToBinary() ([]byte, error) {
	buf := &bytes.Buffer{}
	for _, child := range t.V {
		_, err := child.WriteTo(buf)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
