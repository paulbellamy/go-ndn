package encoding

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Writer(t *testing.T) {
	buf := &bytes.Buffer{}
	w := NewWriter(buf)
	n, err := w.Write(ByteTLV{123, []byte("foo")})
	assert.NoError(t, err)
	assert.Equal(t, n, int(5))
	assert.Equal(t, buf.Bytes(), []byte{
		123,           // type
		3,             // length
		'f', 'o', 'o', // value
	})
}
