package ndn

import (
	"testing"

	"github.com/paulbellamy/go-ndn/encoding/tlv"
	"github.com/paulbellamy/go-ndn/name"
	"github.com/paulbellamy/go-ndn/packets"
	"github.com/stretchr/testify/assert"
)

func testKeyChain(t *testing.T) *KeyChain {
	return NewKeyChain()
}

func Test_KeyChain_Sign_Data(t *testing.T) {
	subject := testKeyChain(t)
	packet := &packets.Data{}
	packet.SetName(name.New(name.Component{"a"}))
	packet.SetContent([]byte("hello world"))

	assert.Nil(t, packet.GetSignature())
	assert.NoError(t, subject.Sign(packet, name.New(name.Component{"certificate"}), tlv.NewEncoder))
	assert.NotNil(t, packet.GetSignature())
}

func Test_KeyChain_Sign_ErrorWritingPacket(t *testing.T) {
	subject := testKeyChain(t)
	packet := &packets.Data{}

	assert.Nil(t, packet.GetSignature())
	assert.Equal(t, subject.Sign(packet, name.New(name.Component{"certificate"}), tlv.NewEncoder), ErrNameRequired)
	assert.Nil(t, packet.GetSignature())
}

func Test_KeyChain_Sign_CertificateNotFound(t *testing.T) {
	subject := testKeyChain(t)
	packet := &packets.Data{}

	assert.Nil(t, packet.GetSignature())
	assert.Equal(t, subject.Sign(packet, name.New(name.Component{"certificate_not_found"}), tlv.NewEncoder), ErrCertificateNotFound)
	assert.Nil(t, packet.GetSignature())
}
