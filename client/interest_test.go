package client

import (
	"bytes"
	"testing"

	"github.com/paulbellamy/go-ndn/encoding"
	"github.com/stretchr/testify/assert"
)

func Test_Interest_WriteTo(t *testing.T) {
	expected, err := encoding.ParentTLV(
		encoding.InterestType,
		encoding.ParentTLV(
			encoding.NameType,
			encoding.ByteTLV(
				encoding.NameComponentType,
				[]byte{'a'},
			),
		),
	).MarshalBinary()
	assert.NoError(t, err)

	buf := &bytes.Buffer{}
	subject := &Interest{
		Name: Name{"a"},
	}
	n, err := subject.WriteTo(buf)
	assert.NoError(t, err)
	assert.Equal(t, n, int64(len(buf.Bytes())))
	assert.Equal(t, buf.Bytes(), expected)
}

func Test_Interest_WriteTo_WithoutAName(t *testing.T) {
	buf := &bytes.Buffer{}

	subject := &Interest{}
	n, err := subject.WriteTo(buf)
	assert.EqualError(t, err, ErrInterestNameRequired.Error())
	assert.Equal(t, n, int64(0))
	assert.Nil(t, buf.Bytes(), "Expected nothing to be written, but data was found.")
}
