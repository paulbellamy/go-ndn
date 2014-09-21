package encoding

import (
	"bytes"
	"io"
)

type ByteTLV struct {
	T uint64
	V []byte
}

func (t ByteTLV) Type() uint64 {
	return t.T
}

func (t ByteTLV) WriteTo(w io.Writer) (n int64, err error) {
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

func (t ByteTLV) MarshalBinary() ([]byte, error) {
	buf := &bytes.Buffer{}
	_, err := t.WriteTo(buf)
	return buf.Bytes(), err
}
