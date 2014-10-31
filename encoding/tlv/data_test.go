package tlv

import (
	"testing"

	"github.com/paulbellamy/go-ndn/encoding"
	"github.com/paulbellamy/go-ndn/name"
	"github.com/paulbellamy/go-ndn/packets"
	"github.com/stretchr/testify/assert"
)

func Test_marshalDataPacket(t *testing.T) {
	packet := &packets.Data{}
	packet.SetName(name.New(name.Component{"a"}))
	packet.SetContent([]byte("hello world"))
	packet.SetSignature(packets.Sha256WithRSASignature([]byte("abcd1234")))

	result, err := marshalDataPacket(packet)
	assert.NoError(t, err)
	assert.Equal(t, result, ParentTLV{
		DataType,
		[]TLV{
			ParentTLV{
				NameType,
				[]TLV{
					ByteTLV{NameComponentType, []byte{'a'}},
				},
			},
			ParentTLV{MetaInfoType, []TLV{}},
			ByteTLV{ContentType, []byte("hello world")},
			ParentTLV{
				SignatureInfoType,
				[]TLV{
					UintTLV{SignatureTypeType, 1 /* Sha256WithRSASignature */},
					// optional KeyLocator
				},
			},
			ByteTLV{SignatureValueType, []byte("abcd1234")},
		},
	})
}

func Test_marshalDataPacket_WithoutAName(t *testing.T) {
	result, err := marshalDataPacket(&packets.Data{})
	assert.Equal(t, err, &encoding.NameRequiredError{})
	assert.Nil(t, result)
}
