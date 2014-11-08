package tlv

import (
	"io"
	"testing"

	"github.com/paulbellamy/go-ndn/encoding"
	"github.com/stretchr/testify/assert"
)

func Test_tlv(t *testing.T) {
	result, rest, err := tlv(1, Bytes).Parse([]byte{0x1, 11, 'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd', '1', '2', '3', '4'})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte("1234"))
	assert.Equal(t, result, GenericTLV{
		T: 1,
		V: []byte("hello world"),
	})
}

func Test_tlv_Small(t *testing.T) {
	result, rest, err := tlv(NameType, zeroOrMore(NameComponent)).Parse([]byte{byte(NameType), 5, byte(NameComponentType), 3, 'f', 'o', 'o'})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{})
	assert.Equal(t, result, GenericTLV{
		T: NameType,
		V: []interface{}{
			GenericTLV{T: NameComponentType, V: []byte("foo")},
		},
	})
}

func Test_tlv_ChecksType(t *testing.T) {
	input := []byte{byte(DataType), 5, byte(NameComponentType), 3, 'f', 'o', 'o'}
	result, rest, err := tlv(NameType, zeroOrMore(NameComponent)).Parse(input)
	assert.EqualError(t, err, "unexpected tlv type")
	assert.Equal(t, rest, input)
	assert.Nil(t, result)
}

func Test_Byte(t *testing.T) {
	result, rest, err := Byte.Parse([]byte{'h'})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{})
	assert.Equal(t, result, byte('h'))
}

func Test_Byte_Short(t *testing.T) {
	result, rest, err := Byte.Parse([]byte{})
	assert.Equal(t, err, io.ErrUnexpectedEOF)
	assert.Equal(t, rest, []byte{})
	assert.Nil(t, result)
}

func Test_Byte_Long(t *testing.T) {
	result, rest, err := Byte.Parse([]byte("hello world"))
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte("ello world"))
	assert.Equal(t, result, byte('h'))
}

func Test_Bytes(t *testing.T) {
	result, rest, err := Bytes.Parse([]byte("hello world"))
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{})
	assert.Equal(t, result, []byte("hello world"))
}

func Test_Bytes_WithLength(t *testing.T) {
	result, rest, err := Bytes(4).Parse([]byte("abcd"))
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{})
	assert.Equal(t, result, []byte("abcd"))
}

func Test_Bytes_WithLength_Short(t *testing.T) {
	result, rest, err := Bytes(4).Parse([]byte("a"))
	assert.Equal(t, err, io.ErrUnexpectedEOF)
	assert.Equal(t, rest, []byte("a"))
	assert.Nil(t, result)
}

func Test_Bytes_WithLength_Long(t *testing.T) {
	result, rest, err := Bytes(4).Parse([]byte("abcd1234"))
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte("1234"))
	assert.Equal(t, result, []byte("abcd"))
}

func Test_Seq(t *testing.T) {
	result, rest, err := seq(Bytes(4), nonNegativeInteger).Parse([]byte("abcd1"))
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{})
	assert.Equal(t, result, []interface{}{[]byte("abcd"), uint64('1')})
}

func Test_Seq_OneParser(t *testing.T) {
	result, rest, err := seq(Bytes(4)).Parse([]byte("abcd1"))
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{'1'})
	assert.Equal(t, result, []byte("abcd"))
}

func Test_Seq_EliminatesNils(t *testing.T) {
	result, rest, err := seq(Bytes(4), maybe(Byte), maybe(Byte)).Parse([]byte("abcd"))
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{})
	assert.Equal(t, result, []interface{}{[]byte("abcd")})
}

func Test_Seq_Short(t *testing.T) {
	result, rest, err := seq(Bytes(1), nonNegativeInteger).Parse([]byte("a"))
	assert.Equal(t, err, io.ErrUnexpectedEOF)
	assert.Equal(t, rest, []byte("a"))
	assert.Nil(t, result)
}

func Test_Seq_Long(t *testing.T) {
	result, rest, err := seq(Bytes(4), Bytes(1)).Parse([]byte("abcd1234"))
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte("234"))
	assert.Equal(t, result, []interface{}{[]byte("abcd"), []byte("1")})
}

func Test_NonNegativeInteger_OneOctet(t *testing.T) {
	result, rest, err := nonNegativeInteger.Parse([]byte{0xff})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{})
	assert.Equal(t, result, uint64(0xff))
}

func Test_NonNegativeInteger_TwoOctet(t *testing.T) {
	result, rest, err := nonNegativeInteger.Parse([]byte{0xff, 0xff})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{})
	assert.Equal(t, result, uint64(0xffff))
}

func Test_NonNegativeInteger_FourOctet(t *testing.T) {
	result, rest, err := nonNegativeInteger.Parse([]byte{0xff, 0xff, 0xff, 0xff})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{})
	assert.Equal(t, result, uint64(0xffffffff))
}

func Test_NonNegativeInteger_EightOctet(t *testing.T) {
	result, rest, err := nonNegativeInteger.Parse([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{})
	assert.Equal(t, result, uint64(0xffffffffffffffff))
}

func Test_NonNegativeInteger_Short(t *testing.T) {
	result, rest, err := nonNegativeInteger.Parse([]byte{})
	assert.Equal(t, err, io.ErrUnexpectedEOF)
	assert.Equal(t, rest, []byte{})
	assert.Nil(t, result)
}

func Test_NonNegativeInteger_Malformed(t *testing.T) {
	result, rest, err := nonNegativeInteger.Parse([]byte{0xff, 2, 3})
	assert.Equal(t, err, io.ErrUnexpectedEOF)
	assert.Equal(t, rest, []byte{0xff, 2, 3})
	assert.Nil(t, result)
}

func Test_Empty(t *testing.T) {
	result, rest, err := empty.Parse([]byte{})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{})
	assert.Nil(t, result)
}

func Test_Empty_Long(t *testing.T) {
	result, rest, err := empty.Parse([]byte{0xff, 2, 3})
	assert.Equal(t, err, &encoding.InvalidUnmarshalError{Message: "unexpected bytes"})
	assert.Equal(t, rest, []byte{0xff, 2, 3})
	assert.Nil(t, result)
}

func Test_Times(t *testing.T) {
	result, rest, err := times(2, Byte).Parse([]byte{0x1, 0x2})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{})
	assert.Equal(t, result, []interface{}{byte(0x1), byte(0x2)})
}

func Test_Times_Short(t *testing.T) {
	result, rest, err := times(2, Byte).Parse([]byte{0x1})
	assert.Equal(t, err, io.ErrUnexpectedEOF)
	assert.Equal(t, rest, []byte{0x1})
	assert.Nil(t, result)
}

func Test_Times_Long(t *testing.T) {
	result, rest, err := times(2, Byte).Parse([]byte{0x1, 0x2, 0x3, 0x4})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{0x3, 0x4})
	assert.Equal(t, result, []interface{}{byte(0x1), byte(0x2)})
}

func Test_ZeroOrMore(t *testing.T) {
	result, rest, err := zeroOrMore(Byte, Byte).Parse([]byte{0x1, 0x2})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{})
	assert.Equal(t, result, []interface{}{[]interface{}{byte(0x1), byte(0x2)}})
}

func Test_ZeroOrMore_Short(t *testing.T) {
	result, rest, err := zeroOrMore(Byte, Byte).Parse([]byte{0x1})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{0x1})
	assert.Equal(t, result, []interface{}{})
}

func Test_ZeroOrMore_Long(t *testing.T) {
	result, rest, err := zeroOrMore(Byte, Byte).Parse([]byte{0x1, 0x2, 0x3, 0x4, 0x5})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{0x5})
	assert.Equal(t, result, []interface{}{[]interface{}{byte(0x1), byte(0x2)}, []interface{}{byte(0x3), byte(0x4)}})
}

func Test_ZeroOrMore_OneParser(t *testing.T) {
	result, rest, err := zeroOrMore(Byte).Parse([]byte{0x1, 0x2})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{})
	assert.Equal(t, result, []interface{}{byte(0x1), byte(0x2)})
}

func Test_OneOrMore(t *testing.T) {
	result, rest, err := oneOrMore(Byte, Byte).Parse([]byte{0x1, 0x2})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{})
	assert.Equal(t, result, []interface{}{[]interface{}{byte(0x1), byte(0x2)}})
}

func Test_OneOrMore_Short(t *testing.T) {
	result, rest, err := oneOrMore(Byte, Byte).Parse([]byte{0x1})
	assert.Equal(t, err, io.ErrUnexpectedEOF)
	assert.Equal(t, rest, []byte{0x1})
	assert.Nil(t, result)
}

func Test_OneOrMore_Long(t *testing.T) {
	result, rest, err := oneOrMore(Byte, Byte).Parse([]byte{0x1, 0x2, 0x3, 0x4, 0x5})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{0x5})
	assert.Equal(t, result, []interface{}{[]interface{}{byte(0x1), byte(0x2)}, []interface{}{byte(0x3), byte(0x4)}})
}

func Test_OneOrMore_OneParser(t *testing.T) {
	result, rest, err := oneOrMore(Byte).Parse([]byte{0x1, 0x2})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{})
	assert.Equal(t, result, []interface{}{byte(0x1), byte(0x2)})
}

func Test_Or(t *testing.T) {
	result, rest, err := or(Bytes(4), Byte).Parse([]byte{0x1, 0x2})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{0x2})
	assert.Equal(t, result, byte(0x1))
}

func Test_Or_NoneMatch(t *testing.T) {
	result, rest, err := or(Bytes(4), Bytes(3)).Parse([]byte{0x1})
	assert.Equal(t, err, io.ErrUnexpectedEOF)
	assert.Equal(t, rest, []byte{0x1})
	assert.Nil(t, result)
}

func Test_Or_Long(t *testing.T) {
	result, rest, err := or(Bytes(4), Bytes(2)).Parse([]byte{0x1, 0x2, 0x3, 0x4, 0x5})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{0x5})
	assert.Equal(t, result, []byte{0x1, 0x2, 0x3, 0x4})
}

func Test_Or_OneParser(t *testing.T) {
	result, rest, err := or(Byte).Parse([]byte{0x1, 0x2})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{0x2})
	assert.Equal(t, result, byte(0x1))
}

func Test_Maybe(t *testing.T) {
	result, rest, err := maybe(Bytes(2)).Parse([]byte{0x1, 0x2})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{})
	assert.Equal(t, result, []byte{0x1, 0x2})
}

func Test_Maybe_NotMatched(t *testing.T) {
	result, rest, err := maybe(Bytes(4)).Parse([]byte{0x1})
	assert.Equal(t, err, nil)
	assert.Equal(t, rest, []byte{0x1})
	assert.Nil(t, result)
}

func Test_Maybe_Long(t *testing.T) {
	result, rest, err := maybe(Bytes(4)).Parse([]byte{0x1, 0x2, 0x3, 0x4, 0x5})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{0x5})
	assert.Equal(t, result, []byte{0x1, 0x2, 0x3, 0x4})
}

func Test_Exactly(t *testing.T) {
	result, rest, err := exactly(0x1, 0x2).Parse([]byte{0x1, 0x2})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{})
	assert.Equal(t, result, []byte{0x1, 0x2})
}

func Test_Exactly_NotMatched(t *testing.T) {
	result, rest, err := exactly(0x2).Parse([]byte{0x1})
	assert.Equal(t, err, io.ErrUnexpectedEOF)
	assert.Equal(t, rest, []byte{0x1})
	assert.Nil(t, result)
}

func Test_Exactly_Short(t *testing.T) {
	result, rest, err := exactly(0x1, 0x2).Parse([]byte{0x1})
	assert.Equal(t, err, io.ErrUnexpectedEOF)
	assert.Equal(t, rest, []byte{0x1})
	assert.Nil(t, result)
}

func Test_Exactly_Long(t *testing.T) {
	result, rest, err := exactly(0x1, 0x2, 0x3, 0x4).Parse([]byte{0x1, 0x2, 0x3, 0x4, 0x5})
	assert.NoError(t, err)
	assert.Equal(t, rest, []byte{0x5})
	assert.Equal(t, result, []byte{0x1, 0x2, 0x3, 0x4})
}
