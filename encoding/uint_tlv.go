package encoding

import (
	"bytes"
	"encoding/binary"
)

func UintTLV(t, v uint64) TLV {
	buf := &bytes.Buffer{}

	if v <= 0xff {
		binary.Write(buf, binary.BigEndian, uint8(v))
	} else if v <= 0xffff {
		binary.Write(buf, binary.BigEndian, uint16(v))
	} else if v <= 0xffffffff {
		binary.Write(buf, binary.BigEndian, uint32(v))
	} else {
		binary.Write(buf, binary.BigEndian, v)
	}
	return ByteTLV(t, buf.Bytes())
}
