package tlv

import (
	"bytes"
	"io"
	"reflect"
	"time"

	"github.com/paulbellamy/go-ndn/encoding"
	"github.com/paulbellamy/go-ndn/name"
	"github.com/paulbellamy/go-ndn/packets"
)

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) encoding.Encoder {
	return &Encoder{
		w: w,
	}
}

func (w *Encoder) Encode(p interface{}) error {
	var t TLV
	var err error
	switch p := p.(type) {
	case Marshaler:
		t, err = p.MarshalTLV()
	case *packets.Interest:
		t, err = marshalInterestPacket(p)
	case *packets.Data:
		t, err = marshalDataPacket(p)
	case TLV:
		t = p
	default:
		err = &encoding.UnsupportedTypeError{Encoding: "tlv", Type: reflect.TypeOf(p)}
	}
	if err != nil {
		return err
	}

	// blagh encoding this then copying is inefficient, but I don't know how to
	// check the maximum size otherwise.
	buf := &bytes.Buffer{}
	_, err = t.WriteTo(buf)
	if err != nil {
		return err
	}

	if len(buf.Bytes()) > encoding.MaxNDNPacketSize {
		return &encoding.PacketTooLargeError{}
	}

	_, err = io.Copy(w.w, buf)
	return err
}

type appender func([]TLV, interface{}) ([]TLV, error)

func marshalChildren(p interface{}, as ...appender) ([]TLV, error) {
	var t []TLV
	var err error
	for _, a := range as {
		t, err = a(t, p)
		if err != nil {
			return t, err
		}
	}
	return t, nil
}

func marshalInterestPacket(p *packets.Interest) (TLV, error) {
	v, err := marshalChildren(
		p,
		appendInterestName,
		appendInterestSelectors,
		appendInterestNonce,
		appendInterestScope,
		appendInterestInterestLifetime,
	)
	if err != nil {
		return nil, err
	}

	return GenericTLV{
		T: InterestType,
		V: v,
	}, nil
}

func appendInterestName(t []TLV, packet interface{}) ([]TLV, error) {
	n := packet.(*packets.Interest).GetName()
	if n.Size() == 0 {
		return nil, &encoding.NameRequiredError{}
	}
	return append(t, marshalName(n)), nil
}

func appendInterestSelectors(t []TLV, packet interface{}) ([]TLV, error) {
	p := packet.(*packets.Interest)

	if m := p.GetMinSuffixComponents(); m != -1 {
		t = append(t, GenericTLV{T: MinSuffixComponentsType, V: uint64(m)})
	}

	if m := p.GetMaxSuffixComponents(); m != -1 {
		t = append(t, GenericTLV{T: MaxSuffixComponentsType, V: uint64(m)})
	}

	// TODO: PublisherPublicKeyLocator here

	if e := p.GetExclude(); e != nil {
		t = append(t, marshalExclude(e))
	}

	if s := p.GetChildSelector(); s != -1 {
		t = append(t, GenericTLV{T: ChildSelectorType, V: uint64(s)})
	}

	if p.GetMustBeFresh() {
		t = append(t, GenericTLV{T: MustBeFreshType, V: []byte{}})
	}

	return t, nil
}

func marshalExclude(e *name.Exclude) TLV {
	v := []TLV{}
	for _, component := range []name.Component(*e) {
		v = append(v, marshalNameComponent(component))
	}

	return GenericTLV{T: ExcludeType, V: v}
}

func appendInterestNonce(t []TLV, packet interface{}) ([]TLV, error) {
	v := []byte{}
	for _, b := range packet.(*packets.Interest).GetNonce() {
		v = append(v, b)
	}
	return append(t, GenericTLV{T: NonceType, V: v}), nil
}

func appendInterestScope(t []TLV, packet interface{}) ([]TLV, error) {
	if s := packet.(*packets.Interest).GetScope(); s == -1 {
		return t, nil
	} else {
		return append(t, GenericTLV{T: ScopeType, V: uint64(s)}), nil
	}
}

func appendInterestInterestLifetime(t []TLV, packet interface{}) ([]TLV, error) {
	if l := packet.(*packets.Interest).GetInterestLifetime(); l == -1 {
		return t, nil
	} else {
		return append(t, GenericTLV{T: InterestLifetimeType, V: uint64(l / time.Millisecond)}), nil
	}
}

func marshalName(n name.Name) TLV {
	v := []TLV{}
	for _, component := range n {
		v = append(v, marshalNameComponent(component))
	}

	return GenericTLV{
		T: NameType,
		V: v,
	}
}

func marshalNameComponent(c name.Component) TLV {
	if c == name.Any {
		return GenericTLV{
			T: AnyType,
			V: []byte{},
		}
	}
	return GenericTLV{
		T: NameComponentType,
		V: c.Bytes(),
	}
}

func marshalDataPacket(p *packets.Data) (TLV, error) {
	v, err := marshalChildren(
		p,
		appendDataName,
		appendDataMetaInfo,
		appendDataContent,
		appendDataSignature,
	)
	if err != nil {
		return nil, err
	}

	return GenericTLV{
		T: DataType,
		V: v,
	}, nil
}

func appendDataName(t []TLV, packet interface{}) ([]TLV, error) {
	n := packet.(*packets.Data).GetName()
	if n.Size() == 0 {
		return nil, &encoding.NameRequiredError{}
	}
	return append(t, marshalName(n)), nil
}

func appendDataMetaInfo(t []TLV, packet interface{}) ([]TLV, error) {
	return append(t, marshalMetaInfo(packet.(*packets.Data).MetaInfo)), nil
}

func marshalMetaInfo(m packets.MetaInfo) TLV {
	v := []TLV{}
	if t := m.GetContentType(); t != packets.UnknownContentType {
		v = append(v, GenericTLV{T: ContentTypeType, V: uint64(t)})
	}

	if p := m.GetFreshnessPeriod(); p != -1 {
		v = append(v, GenericTLV{T: FreshnessPeriodType, V: uint64(p / time.Millisecond)})
	}

	if f := m.GetFinalBlockID(); len(f.GetValue()) > 0 {
		v = append(v, GenericTLV{T: FinalBlockIdType, V: []TLV{marshalNameComponent(f)}})
	}

	return GenericTLV{T: MetaInfoType, V: v}
}

func appendDataContent(t []TLV, packet interface{}) ([]TLV, error) {
	return append(t, GenericTLV{T: ContentType, V: packet.(*packets.Data).GetContent()}), nil
}

func appendDataSignature(t []TLV, packet interface{}) ([]TLV, error) {
	s := packet.(*packets.Data).GetSignature()
	if s == nil {
		return nil, &encoding.SignatureRequiredError{}
	}
	return append(t, marshalSignature(s)...), nil
}

func marshalSignature(s packets.Signature) []TLV {
	return []TLV{
		GenericTLV{
			T: SignatureInfoType,
			V: []TLV{
				GenericTLV{T: SignatureTypeType, V: s.Type()},
				//TODO: GenericTLV{T: KeyLocatorType, V: s.Type()},
			},
		},
		GenericTLV{
			T: SignatureValueType,
			V: s.Bytes(),
		},
	}
}
