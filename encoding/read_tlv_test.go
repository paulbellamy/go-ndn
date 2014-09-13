package encoding

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReadTLV_EOF(t *testing.T) {
	tlv, err := ReadTLV(bytes.NewReader([]byte{}))
	assert.EqualError(t, err, io.EOF.Error())
	assert.Nil(t, tlv)
}

func Test_ReadTLV_UnderflowOnType(t *testing.T) {
	tlv, err := ReadTLV(bytes.NewReader([]byte{255}))
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	assert.Nil(t, tlv)
}

func Test_ReadTLV_UnderflowOnLength(t *testing.T) {
	tlv, err := ReadTLV(bytes.NewReader([]byte{1, 255}))
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	assert.Nil(t, tlv)
}

func Test_ReadTLV_UnderflowOnValue(t *testing.T) {
	tlv, err := ReadTLV(bytes.NewReader([]byte{1, 4, 0}))
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	assert.Nil(t, tlv)
}

func Test_ReadTLV(t *testing.T) {
	tlv, err := ReadTLV(bytes.NewReader([]byte{123, 3, 'f', 'o', 'o'}))
	assert.NoError(t, err)
	assert.Equal(t, tlv, &TLV{Type: 123, Value: []byte("foo")})
}

func Test_TLVReadFrom(t *testing.T) {
	tlv := &TLV{}
	n, err := tlv.ReadFrom(bytes.NewReader([]byte{123, 3, 'f', 'o', 'o'}))
	assert.Equal(t, n, 5)
	assert.NoError(t, err)
	assert.Equal(t, tlv, &TLV{Type: 123, Value: []byte("foo")})
}
