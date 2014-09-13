package encoding

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WriteTLV(t *testing.T) {
	buf := &bytes.Buffer{}
	err := WriteTLV(buf, &TLV{Type: 123, Value: []byte("foo")})
	assert.NoError(t, err)
	assert.Equal(t, buf.Bytes(), []byte{
		123,           // type
		3,             // length
		'f', 'o', 'o', // value
	})
}

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
