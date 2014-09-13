package encoding

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReadUintTLV_EOF(t *testing.T) {
	tlv := &UintTLV{}
	n, err := tlv.ReadFrom(bytes.NewReader([]byte{}))
	assert.EqualError(t, err, io.EOF.Error())
	assert.Equal(t, n, 0)
}

func Test_ReadUintTLV_UnderflowOnType(t *testing.T) {
	tlv := &UintTLV{}
	n, err := tlv.ReadFrom(bytes.NewReader([]byte{255}))
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	assert.Equal(t, n, 1)
	assert.Equal(t, tlv.Type, uint64(0))
}

func Test_ReadUintTLV_IncorrectLengthValue(t *testing.T) {
	tlv := &UintTLV{}
	n, err := tlv.ReadFrom(bytes.NewReader([]byte{1, 3, 0, 0, 0}))
	assert.EqualError(t, err, ErrUnexpectexUintTLVLengthValue.Error())
	assert.Equal(t, n, 2)
	assert.Equal(t, tlv.Type, uint64(1))
}

func Test_ReadUintTLV_UnderflowOnValue(t *testing.T) {
	tlv := &UintTLV{}
	n, err := tlv.ReadFrom(bytes.NewReader([]byte{1, 4, 0}))
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	assert.Equal(t, n, 2)
	assert.Equal(t, tlv.Type, uint64(1))
}

func Test_ReadUintTLV_OneOctetValue(t *testing.T) {
	tlv := &UintTLV{}
	n, err := tlv.ReadFrom(bytes.NewReader([]byte{123, 1, 0xff}))
	assert.NoError(t, err)
	assert.Equal(t, n, 3)
	assert.Equal(t, tlv, &UintTLV{Type: 123, Value: 0xff})
}

func Test_ReadUintTLV_TwoOctetValue(t *testing.T) {
	tlv := &UintTLV{}
	n, err := tlv.ReadFrom(bytes.NewReader([]byte{123, 2, 0xff, 0xff}))
	assert.NoError(t, err)
	assert.Equal(t, n, 4)
	assert.Equal(t, tlv, &UintTLV{Type: 123, Value: 0xffff})
}

func Test_ReadUintTLV_FourOctetValue(t *testing.T) {
	tlv := &UintTLV{}
	n, err := tlv.ReadFrom(bytes.NewReader([]byte{123, 4, 0xff, 0xff, 0xff, 0xff}))
	assert.NoError(t, err)
	assert.Equal(t, n, 6)
	assert.Equal(t, tlv, &UintTLV{Type: 123, Value: 0xffffffff})
}

func Test_ReadUintTLV_EightOctetValue(t *testing.T) {
	tlv := &UintTLV{}
	n, err := tlv.ReadFrom(bytes.NewReader([]byte{
		123,
		8,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
	}))
	assert.NoError(t, err)
	assert.Equal(t, n, 10)
	assert.Equal(t, tlv, &UintTLV{Type: 123, Value: 0xffffffffffffffff})
}
