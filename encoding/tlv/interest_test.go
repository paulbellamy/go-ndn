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

func Test_marshalInterestPacket(t *testing.T) {
	packet := &packets.Interest{}
	packet.SetName(name.New(name.Component{"a"}))

	result, err := marshalInterestPacket(packet)
	assert.NoError(t, err)
	assert.Equal(t, result, ParentTLV{
		InterestType,
		[]TLV{
			ParentTLV{NameType, []TLV{ByteTLV{NameComponentType, []byte{'a'}}}},
		},
	})
}

func Test_marshalInterestPacket_BlankName(t *testing.T) {
	result, err := marshalInterestPacket(&packets.Interest{})
	assert.Equal(t, err, &encoding.NameRequiredError{})
	assert.Nil(t, result)
}

func Test_unmarshalInterestPacket(t *testing.T) {
	var result interface{}

	buffer := &bytes.Buffer{}
	_, err := marshalName(name.New(name.Component{"a"})).WriteTo(buffer)
	assert.NoError(t, err)

	expected := &packets.Interest{}
	expected.SetName(name.New(name.Component{"a"}))

	err = unmarshalInterestPacket(reflect.ValueOf(&result), buffer.Bytes())
	if assert.NoError(t, err) {
		assert.Equal(t, result, expected)
	}
}

func Test_unmarshalInterestPacket_FirstTLVIsNotAName(t *testing.T) {
	t.Error("pending")
}

func Test_unmarshalInterestPacket_InvalidName(t *testing.T) {
	t.Error("pending")
}
