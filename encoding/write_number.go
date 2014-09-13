package encoding

import (
	"encoding/binary"
	"io"
)

func WriteNumber(w io.Writer, num uint64) error {
	if num < 253 {
		return binary.Write(w, binary.BigEndian, uint8(num))
	}

	var prefix uint8
	var value interface{}
	if num <= 0xffff {
		prefix = 253
		value = uint16(num)
	} else if num <= 0xffffffff {
		prefix = 254
		value = uint32(num)
	} else {
		prefix = 255
		value = uint64(num)
	}
	err := binary.Write(w, binary.BigEndian, prefix)
	if err != nil {
		return err
	}

	return binary.Write(w, binary.BigEndian, value)
}
