package tlv

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WriteNumber_OneOctetValue(t *testing.T) {
	buf := &bytes.Buffer{}
	n, err := WriteNumber(buf, 252)
	assert.NoError(t, err)
	assert.Equal(t, n, int64(1))
	assert.Equal(t, buf.Bytes(), []byte{252})
}

func Test_WriteNumber_TwoOctetValue(t *testing.T) {
	buf := &bytes.Buffer{}
	n, err := WriteNumber(buf, 512)
	assert.NoError(t, err)
	assert.Equal(t, n, int64(3))
	assert.Equal(t, buf.Bytes(), []byte{253, 0x02, 00})
}

func Test_WriteNumber_FourOctetValue(t *testing.T) {
	buf := &bytes.Buffer{}
	n, err := WriteNumber(buf, 0xFFFFFFFF)
	assert.NoError(t, err)
	assert.Equal(t, n, int64(5))
	assert.Equal(t, buf.Bytes(), []byte{
		254,
		0xFF, 0xFF, 0xFF, 0xFF})
}

func Test_WriteNumber_EightOctetValue(t *testing.T) {
	buf := &bytes.Buffer{}
	n, err := WriteNumber(buf, 0xFFFFFFFFFFFFFFFF)
	assert.NoError(t, err)
	assert.Equal(t, n, int64(9))
	assert.Equal(t, buf.Bytes(), []byte{
		255,
		0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF})
}
