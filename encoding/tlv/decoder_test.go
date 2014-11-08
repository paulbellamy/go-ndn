package tlv

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/paulbellamy/go-ndn/encoding"
	"github.com/paulbellamy/go-ndn/name"
	"github.com/paulbellamy/go-ndn/packets"
	"github.com/stretchr/testify/assert"
)

/*
func Test_ReadTLV_EOF(t *testing.T) {
	packet, err := NewDecoder(bytes.NewReader([]byte{})).Decode()
	assert.EqualError(t, err, io.EOF.Error())
	assert.Nil(t, packet)
}

func Test_ReadTLV_UnderflowOnType(t *testing.T) {
	packet, err := NewDecoder(bytes.NewReader([]byte{255})).Decode()
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	assert.Nil(t, packet)
}

func Test_ReadTLV_UnderflowOnLength(t *testing.T) {
	packet, err := NewDecoder(bytes.NewReader([]byte{1, 255})).Decode()
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	assert.Nil(t, packet)
}

func Test_ReadTLV_UnderflowOnValue(t *testing.T) {
	packet, err := NewDecoder(bytes.NewReader([]byte{1, 4, 0})).Decode()
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	assert.Nil(t, packet)
}

func Test_ReadTLV(t *testing.T) {
	packet, err := NewDecoder(bytes.NewReader([]byte{123, 3, 'f', 'o', 'o'})).Decode()
	assert.NoError(t, err)
	assert.Equal(t, packet, ByteTLV{T: 123, V: []byte("foo")})
}

func Test_ReadParentTLV(t *testing.T) {
	buf := &bytes.Buffer{}
	expected := ParentTLV{
		T: DataType,
		V: []TLV{
			ParentTLV{
				T: NameType,
				V: []TLV{
					ByteTLV{
						T: NameComponentType,
						V: []byte("a"),
					},
				},
			},
		},
	}
	_, err := expected.WriteTo(buf)
	assert.NoError(t, err)

	packet, err := NewDecoder(buf).Decode()
	assert.NoError(t, err)
	assert.Equal(t, packet, expected)
}

/*
func Test_ReadUintTLV_UnderflowOnType(t *testing.T) {
	err := NewDecoder(bytes.NewReader([]byte{255})).Read()
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	assert.Nil(t, tlv)
}

func Test_ReadUintTLV_IncorrectLengthValue(t *testing.T) {
	err := NewDecoder(bytes.NewReader([]byte{1, 3, 0, 0, 0})).Read()
	assert.NoError(t, err)
	assert.Equal(t, tlv.Type, uint64(1))

	value, err := tlv.Uint()
	assert.Equal(t, value, uint64(0))
	assert.EqualError(t, err, ErrUnexpectexUintTLVLengthValue.Error())
}

func Test_ReadUintTLV_UnderflowOnValue(t *testing.T) {
	err := NewDecoder(bytes.NewReader([]byte{1, 4, 0})).Read()
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	assert.Nil(t, tlv)
}

func Test_ReadUintTLV_OneOctetValue(t *testing.T) {
	err := NewDecoder(bytes.NewReader([]byte{123, 1, 0xff})).Read()
	assert.NoError(t, err)
	assert.Equal(t, tlv, UintTLV(123, 0xff))

	value, err := tlv.Uint()
	assert.NoError(t, err)
	assert.Equal(t, value, uint64(0xff))
}

func Test_ReadUintTLV_TwoOctetValue(t *testing.T) {
	err := NewDecoder(bytes.NewReader([]byte{123, 2, 0xff, 0xff})).Read()
	assert.NoError(t, err)
	assert.Equal(t, tlv, UintTLV(123, 0xffff))

	value, err := tlv.Uint()
	assert.NoError(t, err)
	assert.Equal(t, value, uint64(0xffff))
}

func Test_ReadUintTLV_FourOctetValue(t *testing.T) {
	err := NewDecoder(bytes.NewReader([]byte{123, 4, 0xff, 0xff, 0xff, 0xff})).Read()
	assert.NoError(t, err)
	assert.Equal(t, tlv, UintTLV(123, 0xffffffff))

	value, err := tlv.Uint()
	assert.NoError(t, err)
	assert.Equal(t, value, uint64(0xffffffff))
}

func Test_ReadUintTLV_EightOctetValue(t *testing.T) {
	err := NewDecoder(bytes.NewReader([]byte{
		123,
		8,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
	})).Read()
	assert.NoError(t, err)
	assert.Equal(t, tlv, UintTLV(123, 0xffffffffffffffff))

	value, err := tlv.Uint()
	assert.NoError(t, err)
	assert.Equal(t, value, uint64(0xffffffffffffffff))
}
*/

func Test_Decoder_DataPacket(t *testing.T) {
	buf := bytes.NewReader([]byte{
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

	packet, err := NewDecoder(buf).Decode()
	assert.NoError(t, err)
	if packet == nil {
		t.Fatal("Expected packet not to be nil")
	}

	data, ok := packet.(*packets.Data)
	assert.True(t, ok)
	assert.Equal(t, name.New(name.Component{"foo"}), data.GetName())
	assert.Equal(t, packets.DefaultContentType, data.GetContentType())
	assert.Equal(t, 1*time.Second, data.GetFreshnessPeriod())
	assert.Equal(t, name.Component{"foo"}, data.GetFinalBlockID())
	assert.Equal(t, []byte("hello world"), data.GetContent())
	assert.Equal(t, packets.Sha256WithRSASignature("abcd"), data.GetSignature())
}

func Test_Decoder_DataPacket_Minimal(t *testing.T) {
	buf := bytes.NewReader([]byte{
		byte(DataType), 61,
		byte(NameType), 5,
		byte(NameComponentType), 3, 'f', 'o', 'o',
		byte(MetaInfoType), 0,
		byte(ContentType), 11, 'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd',
		byte(SignatureInfoType), 3,
		byte(SignatureTypeType), 1, 0, /* DigestSha256 */
		byte(SignatureValueType), 32,
		'a', 'b', 'c', 'd', 'a', 'b', 'c', 'd',
		'a', 'b', 'c', 'd', 'a', 'b', 'c', 'd',
		'a', 'b', 'c', 'd', 'a', 'b', 'c', 'd',
		'a', 'b', 'c', 'd', 'a', 'b', 'c', 'd',
	})

	packet, err := NewDecoder(buf).Decode()
	assert.NoError(t, err)
	if packet == nil {
		t.Fatal("Expected packet not to be nil")
	}

	data, ok := packet.(*packets.Data)
	assert.True(t, ok)
	assert.Equal(t, name.New(name.Component{"foo"}), data.GetName())
	assert.Equal(t, packets.UnknownContentType, data.GetContentType())
	assert.Equal(t, -1, data.GetFreshnessPeriod())
	assert.Equal(t, name.Component{}, data.GetFinalBlockID())
	assert.Equal(t, []byte("hello world"), data.GetContent())
	assert.Equal(t, packets.DigestSha256("abcdabcdabcdabcdabcdabcdabcdabcd"), data.GetSignature())
}

func Test_Decoder_DataPacket_UnderflowBytes(t *testing.T) {
	buf := bytes.NewReader([]byte{
		byte(DataType), 2,
		byte(NameType), 5,
	})

	packet, err := NewDecoder(buf).Decode()
	assert.Equal(t, err, io.ErrUnexpectedEOF)
	assert.Nil(t, packet)
}

func Test_Decoder_DataPacket_LeftoverBytes(t *testing.T) {
	buf := bytes.NewReader([]byte{
		byte(DataType), 35,
		byte(NameType), 5,
		byte(NameComponentType), 3, 'f', 'o', 'o',
		byte(MetaInfoType), 0,
		byte(ContentType), 11, 'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd',
		byte(SignatureInfoType), 3,
		byte(SignatureTypeType), 1, 1, /* Sha256WithRSASignature */
		byte(SignatureValueType), 4, 'a', 'b', 'c', 'd',
		0x01, 0x02,
	})

	packet, err := NewDecoder(buf).Decode()
	assert.Equal(t, err, &encoding.InvalidUnmarshalError{Message: "leftover bytes"})
	assert.Nil(t, packet)
}
