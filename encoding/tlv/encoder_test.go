package tlv

import (
	"bytes"
	"reflect"
	"testing"
	"time"

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
	packet.SetMinSuffixComponents(1)
	packet.SetMaxSuffixComponents(1)
	//packet.SetPublisherPublicKeyLocator()
	packet.SetExclude(&name.Exclude{name.Any, name.Component{"a"}})
	packet.SetChildSelector(0)
	packet.SetMustBeFresh(true)
	packet.SetNonce([4]byte{'a', 'b', 'c', 'd'})
	packet.SetScope(1)
	packet.SetInterestLifetime(1 * time.Second)

	err := NewEncoder(buf).Encode(packet)
	assert.NoError(t, err)
	assert.Equal(t, buf.Bytes(), []byte{
		byte(InterestType), 38,
		// Name
		byte(NameType), 5,
		byte(NameComponentType), 3, 'f', 'o', 'o',
		// Selectors?
		//   MinSuffixComponents?
		byte(MinSuffixComponentsType), 1, 1,
		//   MaxSuffixComponents?
		byte(MaxSuffixComponentsType), 1, 1,
		//   PublisherPublicKeyLocator? (not for now)
		//   Exclude?
		byte(ExcludeType), 5, byte(AnyType), 0, byte(NameComponentType), 1, 'a',
		//   ChildSelector?
		byte(ChildSelectorType), 1, 0,
		//   MustBeFresh?
		byte(MustBeFreshType), 0,
		// Nonce
		byte(NonceType), 4, 'a', 'b', 'c', 'd',
		// Scope?
		byte(ScopeType), 1, 1,
		// InterestLifetime? (1000ms)
		byte(InterestLifetimeType), 2, 0x03, 0xe8,
	})
}

func Test_Encoder_InterestPacket_Minimal(t *testing.T) {
	buf := &bytes.Buffer{}
	packet := &packets.Interest{}
	packet.SetName(name.New(name.Component{"foo"}))
	packet.SetMustBeFresh(false)
	packet.SetNonce([4]byte{'a', 'b', 'c', 'd'})

	err := NewEncoder(buf).Encode(packet)
	assert.NoError(t, err)
	assert.Equal(t, buf.Bytes(), []byte{
		byte(InterestType), 13,
		// Name
		byte(NameType), 5,
		byte(NameComponentType), 3, 'f', 'o', 'o',
		// Nonce
		byte(NonceType), 4, 'a', 'b', 'c', 'd',
	})
}

func Test_Encoder_InterestPacket_WithoutAName(t *testing.T) {
	buf := &bytes.Buffer{}
	err := NewEncoder(buf).Encode(&packets.Interest{})
	assert.Equal(t, err, &encoding.NameRequiredError{})
	assert.Nil(t, buf.Bytes(), "Expected nothing to be written, but data was found.")
}

func Test_Encoder_InterestPacket_NonceNotSet(t *testing.T) {
	t.Error("pending, is this even possible to test??")
}

func Test_Encoder_DataPacket(t *testing.T) {
	buf := &bytes.Buffer{}
	packet := &packets.Data{}
	packet.SetName(name.New(name.Component{"foo"}))
	packet.SetContentType(packets.DefaultContentType)
	packet.SetFreshnessPeriod(1 * time.Second)
	packet.SetFinalBlockID(name.Component{"foo"})
	packet.SetContent([]byte("hello world"))
	packet.SetSignature(packets.Sha256WithRSASignature("abcd"))

	err := NewEncoder(buf).Encode(packet)
	assert.NoError(t, err)
	assert.Equal(t, buf.Bytes(), []byte{
		byte(DataType), 47,
		byte(NameType), 5,
		byte(NameComponentType), 3, 'f', 'o', 'o',
		byte(MetaInfoType), 14,
		byte(ContentTypeType), 1, 0,
		byte(FreshnessPeriodType), 2, 0x03, 0xe8,
		byte(FinalBlockIdType), 5, byte(NameComponentType), 3, 'f', 'o', 'o',
		byte(ContentType), 11, 'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd',
		byte(SignatureInfoType), 3,
		byte(SignatureTypeType), 1, 1, /* Sha256WithRSASignature */
		byte(SignatureValueType), 4, 'a', 'b', 'c', 'd',
	})
}

func Test_Encoder_DataPacket_Minimal(t *testing.T) {
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

func Test_Encoder_DataPacket_WithoutAName(t *testing.T) {
	buf := &bytes.Buffer{}
	err := NewEncoder(buf).Encode(&packets.Data{})
	assert.Equal(t, err, &encoding.NameRequiredError{})
	assert.Nil(t, buf.Bytes(), "Expected nothing to be written, but data was found.")
}

func Test_Encoder_DataPacket_WithoutASignature(t *testing.T) {
	buf := &bytes.Buffer{}
	packet := &packets.Data{}
	packet.SetName(name.New(name.Component{"foo"}))
	err := NewEncoder(buf).Encode(packet)
	assert.Equal(t, err, &encoding.SignatureRequiredError{})
	assert.Nil(t, buf.Bytes(), "Expected nothing to be written, but data was found.")
}
