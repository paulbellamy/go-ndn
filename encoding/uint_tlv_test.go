package encoding

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WriteUintTLV_OneOctetValue(t *testing.T) {
	buf := &bytes.Buffer{}
	n, err := UintTLV{123, 0xff}.WriteTo(buf)
	assert.NoError(t, err)
	assert.Equal(t, n, int64(3))
	assert.Equal(t, buf.Bytes(), []byte{
		123,  // type
		1,    // length
		0xff, // value
	})
}

func Test_WriteUintTLV_TwoOctetValue(t *testing.T) {
	buf := &bytes.Buffer{}
	n, err := UintTLV{123, 0xffff}.WriteTo(buf)
	assert.NoError(t, err)
	assert.Equal(t, n, int64(4))
	assert.Equal(t, buf.Bytes(), []byte{
		123,        // type
		2,          // length
		0xff, 0xff, // value
	})
}

func Test_WriteUintTLV_FourOctetValue(t *testing.T) {
	buf := &bytes.Buffer{}
	n, err := UintTLV{123, 0xffffffff}.WriteTo(buf)
	assert.NoError(t, err)
	assert.Equal(t, n, int64(6))
	assert.Equal(t, buf.Bytes(), []byte{
		123,                    // type
		4,                      // length
		0xff, 0xff, 0xff, 0xff, // value
	})
}

func Test_WriteUintTLV_EightOctetValue(t *testing.T) {
	buf := &bytes.Buffer{}
	n, err := UintTLV{123, 0xffffffffffffffff}.WriteTo(buf)
	assert.NoError(t, err)
	assert.Equal(t, n, int64(10))
	assert.Equal(t, buf.Bytes(), []byte{
		123,                    // type
		8,                      // length
		0xff, 0xff, 0xff, 0xff, // value
		0xff, 0xff, 0xff, 0xff,
	})
}
