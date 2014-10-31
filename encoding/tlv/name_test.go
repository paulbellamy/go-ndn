package tlv

import (
	"bytes"
	"testing"

	"github.com/paulbellamy/go-ndn/name"
	"github.com/stretchr/testify/assert"
)

func Test_marshalName(t *testing.T) {
	assert.Equal(t, marshalName(name.New(name.Component{"foo"})), ParentTLV{
		T: NameType,
		V: []TLV{
			ByteTLV{
				T: NameComponentType,
				V: []byte("foo"),
			},
		},
	})
}

/*
func Test_NameFromTLV(t *testing.T) {
	tlv := tlv.ParentTLV{
		T: tlv.NameType,
		V: []tlv.TLV{
			tlv.ByteTLV{
				T: tlv.NameComponentType,
				V: []byte("foo"),
			},
		},
	}
	subject, err := NameFromTLV(tlv)
	assert.NoError(t, err)
	assert.Equal(t, subject, Name{Component{"foo"}})
}

func Test_NameFromTLV_WrongType(t *testing.T) {
	tlv := tlv.ParentTLV{T: tlv.InterestType}
	subject, err := NameFromTLV(tlv)
	assert.EqualError(t, err, "TLV is not a name")
	assert.Nil(t, subject)
}
*/

func Test_unmarshalName(t *testing.T) {
	buffer := &bytes.Buffer{}
	expected := name.New(name.Component{"a"}, name.Component{"b"})
	_, err := marshalName(expected).WriteTo(buffer)
	assert.NoError(t, err)

	n, err := unmarshalName(buffer)
	assert.NoError(t, err)
	assert.Equal(t, n, expected)
}

func Test_unmarshalName_WrongTLVType(t *testing.T) {
	t.Error("pending")
}

func Test_unmarshalName_WrongTLVTypeForComponents(t *testing.T) {
	t.Error("pending")
}

func Test_unmarshalName_UintTLVForComponent(t *testing.T) {
	t.Error("pending")
}
