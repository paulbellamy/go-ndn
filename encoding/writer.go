package encoding

import "io"

type Writer struct {
	w io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w: w,
	}
}

func (w *Writer) Write(t TLV) (int, error) {
	n, err := t.WriteTo(w.w)
	return int(n), err
}
