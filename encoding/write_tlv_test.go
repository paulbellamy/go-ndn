package encoding

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TLVWriteTo(t *testing.T) {
	buf := &bytes.Buffer{}
	err := (&TLV{Type: 123, Value: []byte("foo")}).WriteTo(buf)
	assert.NoError(t, err)
	assert.Equal(t, buf.Bytes(), []byte{
		123,           // type
		3,             // length
		'f', 'o', 'o', // value
	})
}

func Test_TLV_MarshalBinary(t *testing.T) {
	b, err := (&TLV{Type: 123, Value: []byte("foo")}).MarshalBinary()
	assert.NoError(t, err)
	assert.Equal(t, b, []byte{
		123,           // type
		3,             // length
		'f', 'o', 'o', // value
	})
}

func Test_WriteUintTLV_OneOctetValue(t *testing.T) {
	buf := &bytes.Buffer{}
	err := UintTLV(123, 0xff).WriteTo(buf)
	assert.NoError(t, err)
	assert.Equal(t, buf.Bytes(), []byte{
		123,  // type
		1,    // length
		0xff, // value
	})
}

func Test_WriteUintTLV_TwoOctetValue(t *testing.T) {
	buf := &bytes.Buffer{}
	err := UintTLV(123, 0xffff).WriteTo(buf)
	assert.NoError(t, err)
	assert.Equal(t, buf.Bytes(), []byte{
		123,        // type
		2,          // length
		0xff, 0xff, // value
	})
}

func Test_WriteUintTLV_FourOctetValue(t *testing.T) {
	buf := &bytes.Buffer{}
	err := UintTLV(123, 0xffffffff).WriteTo(buf)
	assert.NoError(t, err)
	assert.Equal(t, buf.Bytes(), []byte{
		123,                    // type
		4,                      // length
		0xff, 0xff, 0xff, 0xff, // value
	})
}

func Test_WriteUintTLV_EightOctetValue(t *testing.T) {
	buf := &bytes.Buffer{}
	err := UintTLV(123, 0xffffffffffffffff).WriteTo(buf)
	assert.NoError(t, err)
	assert.Equal(t, buf.Bytes(), []byte{
		123,                    // type
		8,                      // length
		0xff, 0xff, 0xff, 0xff, // value
		0xff, 0xff, 0xff, 0xff,
	})
}
