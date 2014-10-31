package tlv

import (
	"testing"

	"github.com/paulbellamy/go-ndn/packets"
	"github.com/stretchr/testify/assert"
)

func Test_marshalSignature(t *testing.T) {
	signature := packets.Sha256WithRSASignature([]byte("abcd1234"))
	result, err := marshalSignature(signature)
	assert.NoError(t, err)
	assert.Equal(t, result, []TLV{
		ParentTLV{
			SignatureInfoType,
			[]TLV{
				UintTLV{SignatureTypeType, 1 /* SignatureSha256WithRsa */},
				// optional KeyLocator
			},
		},
		ByteTLV{SignatureValueType, []byte("abcd1234")},
	})
}

func Test_marshalSignature_Nil(t *testing.T) {
	result, err := marshalSignature(nil)
	assert.NoError(t, err)
	assert.Equal(t, result, []TLV{})
}
