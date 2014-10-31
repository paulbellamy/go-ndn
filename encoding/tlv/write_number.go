package tlv

import (
	"encoding/binary"
	"io"
)

func WriteNumber(w io.Writer, num uint64) (int64, error) {
	if num < 253 {
		return 1, binary.Write(w, binary.BigEndian, uint8(num))
	}

	var n int64
	var prefix uint8
	var value interface{}
	if num <= 0xffff {
		n = 3
		prefix = 253
		value = uint16(num)
	} else if num <= 0xffffffff {
		prefix = 254
		n = 5
		value = uint32(num)
	} else {
		n = 9
		prefix = 255
		value = uint64(num)
	}
	err := binary.Write(w, binary.BigEndian, prefix)
	if err != nil {
		return 1, err
	}

	return n, binary.Write(w, binary.BigEndian, value)
}
