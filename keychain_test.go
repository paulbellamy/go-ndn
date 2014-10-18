package ndn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testKeyChain(t *testing.T) *KeyChain {
	return NewKeyChain()
}

func Test_KeyChain_Sign_Data(t *testing.T) {
	subject := testKeyChain(t)
	packet := &Data{
		name:    Name{Component{"a"}},
		content: []byte("hello world"),
	}

	assert.Nil(t, packet.GetSignature())
	assert.NoError(t, subject.Sign(packet))
	assert.NotNil(t, packet.GetSignature())
}

func Test_KeyChain_Sign_ErrorWritingPacket(t *testing.T) {
	subject := testKeyChain(t)
	packet := &Data{}

	assert.Nil(t, packet.GetSignature())
	assert.Equal(t, subject.Sign(packet), ErrNameRequired)
	assert.Nil(t, packet.GetSignature())
}
