package encoding

import "io"

type Reader struct {
	r io.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		r: r,
	}
}

func (r *Reader) Read() (*TLV, error) {
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

	return &TLV{Type: t, Value: value[0:n]}, nil
}
