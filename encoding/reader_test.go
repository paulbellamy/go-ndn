package encoding

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReadTLV_EOF(t *testing.T) {
	tlv, err := NewReader(bytes.NewReader([]byte{})).Read()
	assert.EqualError(t, err, io.EOF.Error())
	assert.Nil(t, tlv)
}

func Test_ReadTLV_UnderflowOnType(t *testing.T) {
	tlv, err := NewReader(bytes.NewReader([]byte{255})).Read()
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	assert.Nil(t, tlv)
}

func Test_ReadTLV_UnderflowOnLength(t *testing.T) {
	tlv, err := NewReader(bytes.NewReader([]byte{1, 255})).Read()
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	assert.Nil(t, tlv)
}

func Test_ReadTLV_UnderflowOnValue(t *testing.T) {
	tlv, err := NewReader(bytes.NewReader([]byte{1, 4, 0})).Read()
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	assert.Nil(t, tlv)
}

func Test_ReadTLV(t *testing.T) {
	tlv, err := NewReader(bytes.NewReader([]byte{123, 3, 'f', 'o', 'o'})).Read()
	assert.NoError(t, err)
	assert.Equal(t, tlv, ByteTLV(123, []byte("foo")))
}

/*
func Test_ReadUintTLV_UnderflowOnType(t *testing.T) {
	tlv, err := NewReader(bytes.NewReader([]byte{255})).Read()
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	assert.Nil(t, tlv)
}

func Test_ReadUintTLV_IncorrectLengthValue(t *testing.T) {
	tlv, err := NewReader(bytes.NewReader([]byte{1, 3, 0, 0, 0})).Read()
	assert.NoError(t, err)
	assert.Equal(t, tlv.Type, uint64(1))

	value, err := tlv.Uint()
	assert.Equal(t, value, uint64(0))
	assert.EqualError(t, err, ErrUnexpectexUintTLVLengthValue.Error())
}

func Test_ReadUintTLV_UnderflowOnValue(t *testing.T) {
	tlv, err := NewReader(bytes.NewReader([]byte{1, 4, 0})).Read()
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	assert.Nil(t, tlv)
}

func Test_ReadUintTLV_OneOctetValue(t *testing.T) {
	tlv, err := NewReader(bytes.NewReader([]byte{123, 1, 0xff})).Read()
	assert.NoError(t, err)
	assert.Equal(t, tlv, UintTLV(123, 0xff))

	value, err := tlv.Uint()
	assert.NoError(t, err)
	assert.Equal(t, value, uint64(0xff))
}

func Test_ReadUintTLV_TwoOctetValue(t *testing.T) {
	tlv, err := NewReader(bytes.NewReader([]byte{123, 2, 0xff, 0xff})).Read()
	assert.NoError(t, err)
	assert.Equal(t, tlv, UintTLV(123, 0xffff))

	value, err := tlv.Uint()
	assert.NoError(t, err)
	assert.Equal(t, value, uint64(0xffff))
}

func Test_ReadUintTLV_FourOctetValue(t *testing.T) {
	tlv, err := NewReader(bytes.NewReader([]byte{123, 4, 0xff, 0xff, 0xff, 0xff})).Read()
	assert.NoError(t, err)
	assert.Equal(t, tlv, UintTLV(123, 0xffffffff))

	value, err := tlv.Uint()
	assert.NoError(t, err)
	assert.Equal(t, value, uint64(0xffffffff))
}

func Test_ReadUintTLV_EightOctetValue(t *testing.T) {
	tlv, err := NewReader(bytes.NewReader([]byte{
		123,
		8,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
	})).Read()
	assert.NoError(t, err)
	assert.Equal(t, tlv, UintTLV(123, 0xffffffffffffffff))

	value, err := tlv.Uint()
	assert.NoError(t, err)
	assert.Equal(t, value, uint64(0xffffffffffffffff))
}
*/
