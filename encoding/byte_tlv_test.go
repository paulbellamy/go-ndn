package encoding

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TLVWriteTo(t *testing.T) {
	buf := &bytes.Buffer{}
	n, err := ByteTLV{T: 123, V: []byte("foo")}.WriteTo(buf)
	assert.NoError(t, err)
	assert.Equal(t, n, int64(5))
	assert.Equal(t, buf.Bytes(), []byte{
		123,           // type
		3,             // length
		'f', 'o', 'o', // value
	})
}

func Test_TLV_MarshalBinary(t *testing.T) {
	b, err := ByteTLV{T: 123, V: []byte("foo")}.MarshalBinary()
	assert.NoError(t, err)
	assert.Equal(t, b, []byte{
		123,           // type
		3,             // length
		'f', 'o', 'o', // value
	})
}
