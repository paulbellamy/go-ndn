package tlv

import (
	"bytes"
	"testing"

	"github.com/paulbellamy/go-ndn/encoding"
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

func Test_ReadUint_OneOctetValue(t *testing.T) {
	value, n, err := ReadUint(bytes.NewReader([]byte{0xff}))
	assert.NoError(t, err)
	assert.Equal(t, n, int64(1))
	assert.Equal(t, value, uint64(0xff))
}

func Test_ReadUint_TwoOctetValue(t *testing.T) {
	value, n, err := ReadUint(bytes.NewReader([]byte{0xff, 0xff}))
	assert.NoError(t, err)
	assert.Equal(t, n, int64(2))
	assert.Equal(t, value, uint64(0xffff))
}

func Test_ReadUint_FourOctetValue(t *testing.T) {
	value, n, err := ReadUint(bytes.NewReader([]byte{0xff, 0xff, 0xff, 0xff}))
	assert.NoError(t, err)
	assert.Equal(t, n, int64(4))
	assert.Equal(t, value, uint64(0xffffffff))
}

func Test_ReadUint_EightOctetValue(t *testing.T) {
	value, n, err := ReadUint(bytes.NewReader([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}))
	assert.NoError(t, err)
	assert.Equal(t, n, int64(8))
	assert.Equal(t, value, uint64(0xffffffffffffffff))
}

func Test_ReadUint_Malformed(t *testing.T) {
	value, n, err := ReadUint(bytes.NewReader([]byte{0xff, 0xaa, 'a', 'b', 'c', 'd'}))
	assert.Equal(t, err, &encoding.InvalidUnmarshalError{Message: "malformed uint value"})
	assert.Equal(t, n, int64(6))
	assert.Equal(t, value, uint64(0))
}

func Test_ReadUint_Long(t *testing.T) {
	value, n, err := ReadUint(bytes.NewReader([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 'a', 'b', 'c', 'd'}))
	assert.NoError(t, err)
	assert.Equal(t, n, int64(8))
	assert.Equal(t, value, uint64(0xffffffffffffffff))
}
