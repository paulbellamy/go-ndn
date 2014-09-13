package encoding

import (
	"bytes"
	"io"
)

func (t *TLV) WriteTo(w io.Writer) error {
	err := WriteNumber(w, t.Type)
	if err != nil {
		return err
	}

	err = WriteNumber(w, uint64(len(t.Value)))
	if err != nil {
		return err
	}

	_, err = w.Write(t.Value)
	return err
}

func (t *TLV) MarshalBinary() ([]byte, error) {
	buf := &bytes.Buffer{}
	err := t.WriteTo(buf)
	return buf.Bytes(), err
}
