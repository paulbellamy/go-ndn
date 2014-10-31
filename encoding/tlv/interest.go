package tlv

import (
	"bytes"
	"reflect"

	"github.com/paulbellamy/go-ndn/encoding"
	"github.com/paulbellamy/go-ndn/packets"
)

func marshalInterestPacket(i *packets.Interest) (TLV, error) {
	if i.GetName().IsBlank() {
		return nil, &encoding.NameRequiredError{}
	}
	return ParentTLV{InterestType, []TLV{marshalName(i.GetName())}}, nil
}

func unmarshalInterestPacket(rv reflect.Value, b []byte) error {
	r := bytes.NewReader(b)
	n, err := unmarshalName(r)
	if err != nil {
		return err
	}

	packet := &packets.Interest{}
	packet.SetName(n)
	rv.Set(reflect.ValueOf(packet))

	return nil
}

// Interest ::= INTEREST-TYPE TLV-LENGTH
//                Name
//                Selectors?
//                Nonce
//                Scope?
//                InterestLifetime?
var interest = tlv(InterestType,
	name,
	maybe(selectors),
	nonce,
	maybe(scope),
	maybe(interestLifetime),
)

// Selectors ::= SELECTORS-TYPE TLV-LENGTH
//                 MinSuffixComponents?
//                 MaxSuffixComponents?
//                 PublisherPublicKeyLocator?
//                 Exclude?
//                 ChildSelector?
//                 MustBeFresh?
var selectors = tlv(SelectorsType,
	maybe(minSuffixComponents),
	maybe(maxSuffixComponents),
	maybe(publisherPublicKeyLocator),
	maybe(exclude),
	maybe(childSelector),
	maybe(mustBeFresh),
)

// MinSuffixComponents ::= MIN-SUFFIX-COMPONENTS-TYPE TLV-LENGTH
//                           nonNegativeInteger
var minSuffixComponents = tlv(MinSuffixComponentsType, nonNegativeInteger)

// MaxSuffixComponents ::= MAX-SUFFIX-COMPONENTS-TYPE TLV-LENGTH
//                           nonNegativeInteger
var maxSuffixComponents = tlv(MaxSuffixComponentsType, nonNegativeInteger)

// PublisherPublicKeyLocator ::= KeyLocator
var publisherPublicKeyLocator = keyLocator

// Exclude ::= EXCLUDE-TYPE TLV-LENGTH Any? (NameComponent (Any)?)+
var exclude = tlv(ExcludeType, maybe(any), oneOrMore(nameComponent, maybe(any)))

// Any ::= ANY-TYPE TLV-LENGTH(=0)
var any = tlv(AnyType, empty)

// ChildSelector ::= CHILD-SELECTOR-TYPE TLV-LENGTH
//                     nonNegativeInteger
var childSelector = tlv(ChildSelectorType, nonNegativeInteger)

// MustBeFresh ::= MUST-BE-FRESH-TYPE TLV-LENGTH(=0)
var mustBeFresh = tlv(MustBeFreshType, empty)

// Nonce ::= NONCE-TYPE TLV-LENGTH(=4) BYTE{4}
var nonce = tlv(NonceType, exactBytes(4))

// Guiders

// Scope ::= SCOPE-TYPE TLV-LENGTH nonNegativeInteger
var scope = tlv(ScopeType, nonNegativeInteger)

// InterestLifetime ::= INTEREST-LIFETIME-TYPE TLV-LENGTH nonNegativeInteger
var interestLifetime = tlv(InterestLifetimeType, nonNegativeInteger)
