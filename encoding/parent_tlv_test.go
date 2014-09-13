package encoding

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParentTLV_WriteTo(t *testing.T) {
	subject := ParentTLV(
		NameType,
		ByteTLV(NameComponentType, []byte("abcd")),
		ByteTLV(NameComponentType, []byte("1234")),
	)

	buf := &bytes.Buffer{}
	n, err := subject.WriteTo(buf)
	assert.NoError(t, err)
	assert.Equal(t, n, int64(len(buf.Bytes())))
	assert.Equal(t, buf.Bytes(), []byte{
		// NameType
		0x07,
		// Length
		12,

		// First Component
		// NameComponentType
		0x08,
		// Length
		4,
		// Value
		'a', 'b', 'c', 'd',

		// Second Component
		// NameComponentType
		0x08,
		// Length
		4,
		// Value
		'1', '2', '3', '4',
	})
}
