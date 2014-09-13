package encoding

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Writer(t *testing.T) {
	buf := &bytes.Buffer{}
	w := NewWriter(buf)
	err := w.Write(&TLV{Type: 123, Value: []byte("foo")})
	assert.NoError(t, err)
	assert.Equal(t, buf.Bytes(), []byte{
		123,           // type
		3,             // length
		'f', 'o', 'o', // value
	})
}
