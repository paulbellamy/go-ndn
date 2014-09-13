package encoding

import (
	"bytes"
	"io"
)

func ByteTLV(t uint64, v []byte) TLV {
	return &byteTLV{
		T: t,
		V: v,
	}
}

type byteTLV struct {
	T uint64
	V []byte
}

func (t *byteTLV) Type() uint64 {
	return t.T
}

func (t *byteTLV) WriteTo(w io.Writer) (n int64, err error) {
	n, err = WriteNumber(w, t.T)
	if err != nil {
		return
	}

	written, err := WriteNumber(w, uint64(len(t.V)))
	n += written
	if err != nil {
		return
	}

	written2, err := w.Write(t.V)
	n += int64(written2)
	return
}

func (t *byteTLV) MarshalBinary() ([]byte, error) {
	buf := &bytes.Buffer{}
	_, err := t.WriteTo(buf)
	return buf.Bytes(), err
}
