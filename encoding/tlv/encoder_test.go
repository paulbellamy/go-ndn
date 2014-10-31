package tlv

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/paulbellamy/go-ndn/encoding"
	"github.com/paulbellamy/go-ndn/name"
	"github.com/paulbellamy/go-ndn/packets"
	"github.com/stretchr/testify/assert"
)

func Test_Encoder_UnsupportedType(t *testing.T) {
	buf := &bytes.Buffer{}
	w := NewEncoder(buf)
	err := w.Encode("foobar")
	assert.Equal(t, err, &encoding.UnsupportedTypeError{Encoding: "tlv", Type: reflect.TypeOf("")})
}

func Test_Encoder_MaxPacketSize(t *testing.T) {
	content := []byte{}
	for i := 0; i <= encoding.MaxNDNPacketSize; i++ {
		content = append(content, '0')
	}

	packet := &packets.Data{}
	packet.SetName(name.New(name.Component{"a"}))
	packet.SetContent(content)

	buf := &bytes.Buffer{}
	err := NewEncoder(buf).Encode(packet)
	assert.Equal(t, err, &encoding.PacketTooLargeError{})

	// Check no data was written
	assert.Equal(t, len(buf.Bytes()), 0)
}

func Test_Encoder_InterestPacket(t *testing.T) {
	buf := &bytes.Buffer{}
	packet := &packets.Interest{}
	packet.SetName(name.New(name.Component{"foo"}))

	err := NewEncoder(buf).Encode(packet)
	assert.NoError(t, err)
	assert.Equal(t, buf.Bytes(), []byte{
		byte(InterestType), 7,
		byte(NameType), 5,
		byte(NameComponentType), 3, 'f', 'o', 'o',
	})
}

func Test_Encoder_InterestPacket_WithoutAName(t *testing.T) {
	buf := &bytes.Buffer{}
	err := NewEncoder(buf).Encode(&packets.Interest{})
	assert.Equal(t, err, &encoding.NameRequiredError{})
	assert.Nil(t, buf.Bytes(), "Expected nothing to be written, but data was found.")
}

func Test_Encoder_DataPacket(t *testing.T) {
	buf := &bytes.Buffer{}
	packet := &packets.Data{}
	packet.SetName(name.New(name.Component{"foo"}))
	packet.SetContent([]byte("hello world"))
	packet.SetSignature(packets.Sha256WithRSASignature("abcd"))

	err := NewEncoder(buf).Encode(packet)
	assert.NoError(t, err)
	assert.Equal(t, buf.Bytes(), []byte{
		byte(DataType), 33,
		byte(NameType), 5,
		byte(NameComponentType), 3, 'f', 'o', 'o',
		byte(MetaInfoType), 0,
		byte(ContentType), 11, 'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd',
		byte(SignatureInfoType), 3,
		byte(SignatureTypeType), 1, 1, /* Sha256WithRSASignature */
		byte(SignatureValueType), 4, 'a', 'b', 'c', 'd',
	})
}
