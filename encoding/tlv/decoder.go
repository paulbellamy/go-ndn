package tlv

import (
	"fmt"
	"io"
	"reflect"

	"github.com/paulbellamy/go-ndn/encoding"
)

type Decoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) encoding.Decoder {
	return &Decoder{
		r: r,
	}
}

func (r *Decoder) Decode(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &encoding.InvalidUnmarshalTargetError{Encoding: "tlv", Type: reflect.TypeOf(v)}
	}

	rawTLV, err := readByteTLV(r.r)
	if err != nil {
		return err
	}

	switch rawTLV.Type() {
	case InterestType:
		return unmarshalInterestPacket(rv, rawTLV.V)
	case DataType:
		return unmarshalDataPacket(rv, rawTLV.V)
	default:
		return &encoding.InvalidUnmarshalError{fmt.Sprintf("tlv: unexpected tlv type %d", rawTLV.Type())}
	}
}

func (r *Decoder) isParentTLV(t uint64) bool {
	switch t {
	case InterestType, DataType, NameType, SignatureInfoType:
		return true
	default:
		return false
	}
}
