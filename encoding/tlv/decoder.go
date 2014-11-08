package tlv

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/paulbellamy/go-ndn/encoding"
	"github.com/paulbellamy/go-ndn/name"
	"github.com/paulbellamy/go-ndn/packets"
)

type Decoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) encoding.Decoder {
	return &Decoder{
		r: r,
	}
}

func (r *Decoder) Decode() (interface{}, error) {
	rawTLV, err := readByteTLV(r.r)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	_, err = rawTLV.WriteTo(buf)
	if err != nil {
		return nil, err
	}

	switch rawTLV.Type() {
	case InterestType:
		return unmarshalInterestPacket(buf.Bytes())
	case DataType:
		return unmarshalDataPacket(buf.Bytes())
	default:
		return nil, &encoding.InvalidUnmarshalError{fmt.Sprintf("tlv: unexpected tlv type %d", rawTLV.Type())}
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

func unmarshalInterestPacket(input []byte) (interface{}, error) {
	t, rest, err := Interest.Parse(input)
	if err != nil {
		return nil, err
	}
	if len(rest) > 0 {
		return nil, &encoding.InvalidUnmarshalError{Message: "leftover bytes"}
	}
	packet := &packets.Interest{}
	tlv, ok := t.(GenericTLV)
	if !ok {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected interest tlv"}
	}
	tlvs, ok := tlv.V.([]interface{})
	if !ok {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected interest tlv"}
	}
	err = unmarshalTLVs(
		packet, tlvs,
		unmarshalInterestName,
		unmarshalInterestSelectors,
		unmarshalInterestNonce,
		unmarshalInterestScope,
		unmarshalInterestInterestLifetime,
	)
	if err != nil {
		return nil, err
	}
	return packet, nil
}

func unmarshalInterestName(t []interface{}, packet interface{}) ([]interface{}, error) {
	if len(t) < 1 {
		return nil, io.ErrUnexpectedEOF
	}
	n, err := unmarshalName(t[0])
	if err != nil {
		return nil, err
	}
	packet.(*packets.Interest).SetName(n)
	return t[1:], nil
}

func unmarshalInterestSelectors(t []interface{}, packet interface{}) ([]interface{}, error) {
	if len(t) < 1 {
		return t, nil
	}
	tlv, ok := t[0].(GenericTLV)
	if !ok {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected tlv"}
	}
	if tlv.T != SelectorsType {
		return t, nil
	}
	content, ok := tlv.V.([]interface{})
	if !ok {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected selectors tlv"}
	}
	for _, c := range content {
		err := unmarshalInterestSelector(c, packet.(*packets.Interest))
		if err != nil {
			return nil, err
		}
	}
	return t[1:], nil
}

func unmarshalInterestSelector(t interface{}, packet *packets.Interest) error {
	tlv, ok := t.(GenericTLV)
	if !ok {
		return &encoding.InvalidUnmarshalError{Message: "expected selector field tlv"}
	}
	switch tlv.T {
	case MinSuffixComponentsType:
		v, ok := tlv.V.(uint64)
		if !ok {
			return &encoding.InvalidUnmarshalError{Message: "expected minimum suffix components tlv"}
		}
		packet.SetMinSuffixComponents(int(v))
	case MaxSuffixComponentsType:
		v, ok := tlv.V.(uint64)
		if !ok {
			return &encoding.InvalidUnmarshalError{Message: "expected maximum suffix components tlv"}
		}
		packet.SetMaxSuffixComponents(int(v))
	/*
		TODO: case PublisherPublicKeyLocator:
	*/
	case ExcludeType:
		v, ok := tlv.V.([]interface{})
		if !ok {
			return &encoding.InvalidUnmarshalError{Message: "expected exclude tlv"}
		}
		e, err := unmarshalExcludeComponents(v)
		if err != nil {
			return err
		}
		packet.SetExclude(e)
	case ChildSelectorType:
		v, ok := tlv.V.(uint64)
		if !ok {
			return &encoding.InvalidUnmarshalError{Message: "expected child selector tlv"}
		}
		packet.SetChildSelector(int(v))
	case MustBeFreshType:
		packet.SetMustBeFresh(true)
	default:
		return &encoding.InvalidUnmarshalError{Message: "unknown metainfo field type"}
	}
	return nil
}

func unmarshalExcludeComponents(t []interface{}) (*name.Exclude, error) {
	e := name.Exclude{}
	for _, c := range t {
		switch c := c.(type) {
		case GenericTLV:
			component, err := unmarshalNameComponent(c)
			if err != nil {
				return nil, err
			}
			e = append(e, component)
		case []interface{}:
			components, err := unmarshalExcludeComponents(c)
			if err != nil {
				return nil, err
			}
			e = append(e, (*components)...)
		default:
			return nil, &encoding.InvalidUnmarshalError{Message: "expected exclude component tlv"}
		}
	}
	return &e, nil
}

func unmarshalInterestNonce(t []interface{}, packet interface{}) ([]interface{}, error) {
	if len(t) < 1 {
		return nil, io.ErrUnexpectedEOF
	}
	tlv, ok := t[0].(GenericTLV)
	if !ok || tlv.T != NonceType {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected nonce tlv"}
	}
	v, ok := tlv.V.([]byte)
	if !ok {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected nonce tlv"}
	}
	nonce := [4]byte{}
	for i, _ := range nonce {
		nonce[i] = v[i]
	}
	packet.(*packets.Interest).SetNonce(nonce)
	return t[1:], nil
}

func unmarshalInterestScope(t []interface{}, packet interface{}) ([]interface{}, error) {
	if len(t) < 1 {
		return t, nil
	}
	tlv, ok := t[0].(GenericTLV)
	if !ok || tlv.T != ScopeType {
		return t, nil
	}
	v, ok := tlv.V.(uint64)
	if !ok {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected scope tlv"}
	}
	packet.(*packets.Interest).SetScope(int(v))
	return t[1:], nil
}

func unmarshalInterestInterestLifetime(t []interface{}, packet interface{}) ([]interface{}, error) {
	if len(t) < 1 {
		return t, nil
	}
	tlv, ok := t[0].(GenericTLV)
	if !ok || tlv.T != InterestLifetimeType {
		return t, nil
	}
	v, ok := tlv.V.(uint64)
	if !ok {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected interest lifetime tlv"}
	}
	packet.(*packets.Interest).SetInterestLifetime(time.Duration(v) * time.Millisecond)
	return t[1:], nil
}

func unmarshalDataPacket(input []byte) (interface{}, error) {
	t, rest, err := Data.Parse(input)
	if err != nil {
		return nil, err
	}
	if len(rest) > 0 {
		return nil, &encoding.InvalidUnmarshalError{Message: "leftover bytes"}
	}
	packet := &packets.Data{}
	tlv, ok := t.(GenericTLV)
	if !ok {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected data tlv"}
	}
	tlvs, ok := tlv.V.([]interface{})
	if !ok {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected data tlv"}
	}
	err = unmarshalTLVs(
		packet, tlvs,
		unmarshalDataName,
		unmarshalDataMetaInfo,
		unmarshalDataContent,
		unmarshalDataSignature,
	)
	if err != nil {
		return nil, err
	}
	return packet, nil
}

type unmarshaler func(t []interface{}, packet interface{}) ([]interface{}, error)

func unmarshalTLVs(packet interface{}, t []interface{}, us ...unmarshaler) error {
	var err error
	for _, u := range us {
		t, err = u(t, packet)
		if err != nil {
			return err
		}
	}
	return nil
}

func unmarshalDataName(t []interface{}, packet interface{}) ([]interface{}, error) {
	if len(t) < 1 {
		return nil, io.ErrUnexpectedEOF
	}
	n, err := unmarshalName(t[0])
	if err != nil {
		return nil, err
	}
	packet.(*packets.Data).SetName(n)
	return t[1:], nil
}

func unmarshalName(t interface{}) (name.Name, error) {
	tlv, ok := t.(GenericTLV)
	if !ok || tlv.T != NameType {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected name tlv"}
	}
	componentTLVs, ok := tlv.V.([]interface{})
	if !ok {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected name tlv"}
	}
	components := []name.Component{}
	for _, componentTLV := range componentTLVs {
		c, err := unmarshalNameComponent(componentTLV)
		if err != nil {
			return nil, err
		}
		components = append(components, c)
	}
	return name.New(components...), nil
}

func unmarshalNameComponent(t interface{}) (name.Component, error) {
	c := name.Component{}
	tlv, ok := t.(GenericTLV)
	if !ok {
		return c, &encoding.InvalidUnmarshalError{Message: "expected name component tlv"}
	}
	if tlv.T == AnyType {
		return name.Any, nil
	}
	if tlv.T != NameComponentType {
		return c, &encoding.InvalidUnmarshalError{Message: "expected name component tlv"}
	}
	bs, ok := tlv.V.([]byte)
	if !ok {
		return c, &encoding.InvalidUnmarshalError{Message: "expected name component tlv"}
	}
	return name.ComponentFromBytes(bs), nil
}

func unmarshalDataMetaInfo(t []interface{}, packet interface{}) ([]interface{}, error) {
	if len(t) < 1 {
		return nil, io.ErrUnexpectedEOF
	}
	tlv, ok := t[0].(GenericTLV)
	if !ok || tlv.T != MetaInfoType {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected metainfo tlv"}
	}
	content, ok := tlv.V.([]interface{})
	if !ok {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected metainfo tlv"}
	}
	for _, c := range content {
		err := unmarshalMetaInfoField(c, packet.(*packets.Data))
		if err != nil {
			return nil, err
		}
	}
	return t[1:], nil
}

func unmarshalMetaInfoField(t interface{}, packet *packets.Data) error {
	tlv, ok := t.(GenericTLV)
	if !ok {
		return &encoding.InvalidUnmarshalError{Message: "expected metainfo field tlv"}
	}
	switch tlv.T {
	case ContentTypeType:
		v, ok := tlv.V.(uint64)
		if !ok {
			return &encoding.InvalidUnmarshalError{Message: "expected metainfo content type tlv"}
		}
		packet.SetContentType(packets.ContentType(v))
	case FreshnessPeriodType:
		v, ok := tlv.V.(uint64)
		if !ok {
			return &encoding.InvalidUnmarshalError{Message: "expected metainfo freshness period tlv"}
		}
		packet.SetFreshnessPeriod(time.Duration(v) * time.Millisecond)
	case FinalBlockIdType:
		v, err := unmarshalNameComponent(tlv.V)
		if err != nil {
			return err
		}
		packet.SetFinalBlockID(v)
	default:
		return &encoding.InvalidUnmarshalError{Message: "unknown metainfo field type"}
	}
	return nil
}

func unmarshalDataContent(t []interface{}, packet interface{}) ([]interface{}, error) {
	if len(t) < 1 {
		return nil, io.ErrUnexpectedEOF
	}
	tlv, ok := t[0].(GenericTLV)
	if !ok || tlv.T != ContentType {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected content tlv"}
	}
	content, ok := tlv.V.([]byte)
	if !ok {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected content tlv"}
	}
	packet.(*packets.Data).SetContent(content)
	return t[1:], nil
}

func unmarshalDataSignature(t []interface{}, packet interface{}) ([]interface{}, error) {
	if len(t) < 1 {
		return nil, io.ErrUnexpectedEOF
	}
	sig, ok := t[0].([]interface{})
	if !ok {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected signature type tlv"}
	}
	if len(sig) < 2 {
		return nil, io.ErrUnexpectedEOF
	}
	s, err := unmarshalSignature(sig[0], sig[1])
	if err != nil {
		return nil, err
	}
	packet.(*packets.Data).SetSignature(s)
	return t[1:], nil
}

func unmarshalSignature(info, value interface{}) (packets.Signature, error) {
	v, err := unmarshalSignatureValue(value)
	if err != nil {
		return nil, err
	}
	t, err := unmarshalSignatureInfo(info, v)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// TODO: This should do the keylocator as well
func unmarshalSignatureInfo(info interface{}, value []byte) (packets.Signature, error) {
	tlv, ok := info.(GenericTLV)
	if !ok {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected signature info tlv"}
	}
	if tlv.T != SignatureInfoType {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected signature info tlv"}
	}
	content, ok := tlv.V.(GenericTLV)
	if !ok {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected signature type tlv"}
	}
	return unmarshalSignatureType(content, value)
}

func unmarshalSignatureType(t interface{}, value []byte) (packets.Signature, error) {
	tlv, ok := t.(GenericTLV)
	if !ok {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected signature type tlv"}
	}
	if tlv.T != SignatureTypeType {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected signature type tlv"}
	}
	content, ok := tlv.V.([]byte)
	if !ok {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected signature type tlv"}
	}
	if len(content) != 1 {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected signature type tlv"}
	}

	switch content[0] {
	case 0:
		return packets.DigestSha256(value), nil
	case 1:
		return packets.Sha256WithRSASignature(value), nil
	case 3:
		return packets.Sha256WithECDSASignature(value), nil
	default:
		return nil, &encoding.InvalidUnmarshalError{Message: "unknown signature type"}
	}
}

func unmarshalSignatureValue(value interface{}) ([]byte, error) {
	tlv, ok := value.(GenericTLV)
	if !ok {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected signature value tlv"}
	}
	if tlv.T != SignatureValueType {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected signature value tlv"}
	}
	bs, ok := tlv.V.([]byte)
	if !ok {
		return nil, &encoding.InvalidUnmarshalError{Message: "expected signature value tlv"}
	}
	return bs, nil
}
