package encoding

import "io"

func ReadTLV(r io.Reader) (*TLV, error) {
	t, _, err := ReadNumber(r)
	if err != nil {
		return nil, err
	}

	length, _, err := ReadNumber(r)
	if err != nil {
		return nil, err
	}

	value := make([]byte, length)
	n, err := r.Read(value)
	if err != nil {
		return nil, err
	}

	if uint64(n) < length {
		return nil, io.ErrUnexpectedEOF
	}

	return &TLV{Type: t, Value: value[0:n]}, nil
}

func (t *TLV) ReadFrom(r io.Reader) (int64, error) {
	return 0, nil
}
