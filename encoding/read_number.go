package encoding

import (
	"encoding/binary"
	"io"
)

func ReadNumber(r io.Reader) (uint64, int64, error) {
	var first uint8
	err := binary.Read(r, binary.BigEndian, &first)
	if err != nil {
		return 0, 0, err
	}

	if first < 253 {
		return uint64(first), 1, nil
	}

	if first == 253 {
		var value uint16

		err = binary.Read(r, binary.BigEndian, &value)
		if err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			return 0, 1, err
		}

		return uint64(value), 3, nil
	} else if first == 254 {
		var value uint32

		err = binary.Read(r, binary.BigEndian, &value)
		if err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			return 0, 1, err
		}

		return uint64(value), 5, nil
	}

	var value uint64

	err = binary.Read(r, binary.BigEndian, &value)
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return 0, 1, err
	}

	return uint64(value), 9, nil
}
