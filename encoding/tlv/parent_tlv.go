package tlv

import (
	"bytes"
	"io"
)

type ParentTLV struct {
	T uint64
	V []TLV
}

func (t ParentTLV) Type() uint64 {
	return t.T
}

func (t ParentTLV) WriteTo(w io.Writer) (n int64, err error) {
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

func (t ParentTLV) MarshalBinary() ([]byte, error) {
	buf := &bytes.Buffer{}
	_, err := t.WriteTo(buf)
	return buf.Bytes(), err
}

func (t ParentTLV) marshalChildrenToBinary() ([]byte, error) {
	buf := &bytes.Buffer{}
	for _, child := range t.V {
		_, err := child.WriteTo(buf)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func readParentTLV(r io.Reader) (ParentTLV, error) {
	return ParentTLV{}, nil
}
