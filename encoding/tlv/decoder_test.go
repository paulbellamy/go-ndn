package tlv

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/paulbellamy/go-ndn/encoding"
	"github.com/stretchr/testify/assert"
)

func Test_Decoder_Decode_NonPointer(t *testing.T) {
	var result string

	assert.Equal(
		t,
		NewDecoder(nil).Decode(result),
		&encoding.InvalidUnmarshalTargetError{Encoding: "tlv", Type: reflect.TypeOf(result)},
	)
}

func Test_Decoder_Decode_NilTarget(t *testing.T) {
	var result *string

	assert.Equal(
		t,
		NewDecoder(nil).Decode(result),
		&encoding.InvalidUnmarshalTargetError{Encoding: "tlv", Type: reflect.TypeOf(result)},
	)
}

func Test_ReadTLV_EOF(t *testing.T) {
	var packet interface{}
	err := NewDecoder(bytes.NewReader([]byte{})).Decode(&packet)
	assert.EqualError(t, err, io.EOF.Error())
	assert.Nil(t, packet)
}

func Test_ReadTLV_UnderflowOnType(t *testing.T) {
	var packet interface{}
	err := NewDecoder(bytes.NewReader([]byte{255})).Decode(&packet)
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	assert.Nil(t, packet)
}

func Test_ReadTLV_UnderflowOnLength(t *testing.T) {
	var packet interface{}
	err := NewDecoder(bytes.NewReader([]byte{1, 255})).Decode(&packet)
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	assert.Nil(t, packet)
}

func Test_ReadTLV_UnderflowOnValue(t *testing.T) {
	var packet interface{}
	err := NewDecoder(bytes.NewReader([]byte{1, 4, 0})).Decode(&packet)
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	assert.Nil(t, packet)
}

func Test_ReadTLV(t *testing.T) {
	var packet interface{}
	err := NewDecoder(bytes.NewReader([]byte{123, 3, 'f', 'o', 'o'})).Decode(&packet)
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

	var packet interface{}
	err = NewDecoder(buf).Decode(&packet)
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
