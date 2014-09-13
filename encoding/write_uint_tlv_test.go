package encoding

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WriteUintTLV_OneOctetValue(t *testing.T) {
	buf := &bytes.Buffer{}
	err := WriteUintTLV(buf, &UintTLV{Type: 123, Value: 0xff})
	assert.NoError(t, err)
	assert.Equal(t, buf.Bytes(), []byte{
		123,  // type
		1,    // length
		0xff, // value
	})
}

func Test_WriteUintTLV_TwoOctetValue(t *testing.T) {
	buf := &bytes.Buffer{}
	err := WriteUintTLV(buf, &UintTLV{Type: 123, Value: 0xffff})
	assert.NoError(t, err)
	assert.Equal(t, buf.Bytes(), []byte{
		123,        // type
		2,          // length
		0xff, 0xff, // value
	})
}

func Test_WriteUintTLV_FourOctetValue(t *testing.T) {
	buf := &bytes.Buffer{}
	err := WriteUintTLV(buf, &UintTLV{Type: 123, Value: 0xffffffff})
	assert.NoError(t, err)
	assert.Equal(t, buf.Bytes(), []byte{
		123,                    // type
		4,                      // length
		0xff, 0xff, 0xff, 0xff, // value
	})
}

func Test_WriteUintTLV_EightOctetValue(t *testing.T) {
	buf := &bytes.Buffer{}
	err := WriteUintTLV(buf, &UintTLV{Type: 123, Value: 0xffffffffffffffff})
	assert.NoError(t, err)
	assert.Equal(t, buf.Bytes(), []byte{
		123,                    // type
		8,                      // length
		0xff, 0xff, 0xff, 0xff, // value
		0xff, 0xff, 0xff, 0xff,
	})
}

func Test_UintTLVWriteTo_OneOctetValue(t *testing.T) {
	buf := &bytes.Buffer{}
	err := (&UintTLV{Type: 123, Value: 0xff}).WriteTo(buf)
	assert.NoError(t, err)
	assert.Equal(t, buf.Bytes(), []byte{
		123,  // type
		1,    // length
		0xff, // value
	})
}
