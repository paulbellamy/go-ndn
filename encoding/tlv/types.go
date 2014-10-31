package tlv

const (
	// Packet Types
	InterestType uint64 = 0x05
	DataType     uint64 = 0x06

	// Common Fields
	NameType          uint64 = 0x07
	NameComponentType uint64 = 0x08

	// Interest Packet
	SelectorsType        uint64 = 0x09
	NonceType            uint64 = 0x0A
	ScopeType            uint64 = 0x0B
	InterestLifetimeType uint64 = 0x0C

	// Interest/Selectors
	MinSuffixComponentsType       uint64 = 0x0D
	MaxSuffixComponentsType       uint64 = 0x0E
	PublisherPublicKeyLocatorType uint64 = 0x0F
	ExcludeType                   uint64 = 0x10
	ChildSelectorType             uint64 = 0x11
	MustBeFreshType               uint64 = 0x12
	AnyType                       uint64 = 0x13

	// Data Packet
	MetaInfoType       uint64 = 0x14
	ContentType        uint64 = 0x15
	SignatureInfoType  uint64 = 0x16
	SignatureValueType uint64 = 0x17

	// Data/MetaInfo
	ContentTypeType     uint64 = 0x18
	FreshnessPeriodType uint64 = 0x19
	FinalBlockIdType    uint64 = 0x1A

	// Data/Signature
	SignatureTypeType uint64 = 0x1B
	KeyLocatorType    uint64 = 0x1C
	KeyDigestType     uint64 = 0x1D
)

var tlvTypeFactories = map[uint64]string{
	// Packet Types
	InterestType: "ParentTLV",
	DataType:     "ParentTLV",

	// Common Fields
	NameType:          "ParentTLV",
	NameComponentType: "ByteTLV",

	// Interest Packet
	SelectorsType:        "ParentTLV",
	NonceType:            "ByteTLV",
	ScopeType:            "UintTLV",
	InterestLifetimeType: "UintTLV",

	// Interest/Selectors
	MinSuffixComponentsType:       "UintTLV",
	MaxSuffixComponentsType:       "UintTLV",
	PublisherPublicKeyLocatorType: "ParentTLV",
	ExcludeType:                   "ParentTLV",
	ChildSelectorType:             "UintTLV",
	MustBeFreshType:               "ByteTLV",
	AnyType:                       "ByteTLV",

	// Data Packet
	MetaInfoType:       "ParentTLV",
	ContentType:        "ByteTLV",
	SignatureInfoType:  "ParentTLV",
	SignatureValueType: "ByteTLV",

	// Data/MetaInfo
	ContentTypeType:     "UintTLV",
	FreshnessPeriodType: "UintTLV",
	FinalBlockIdType:    "ParentTLV",

	// Data/Signature
	SignatureTypeType: "UintTLV",
	KeyLocatorType:    "ParentTLV",
	KeyDigestType:     "ByteTLV",
}

type parser interface{}

var nonNegativeInteger parser
var bytes parser
var empty = exactBytes(0)

func exactBytes(length int) parser {
	return parser{}
}

func tlv(t uint64, p ...parser) parser {
	return p
}

// zero-or-more sequences repeated
func many(p ...parser) parser {
	return p
}

// one-or-more sequences repeated
func oneOrMore(p ...parser) parser {
	return p
}

// first parser to match
func or(p parser, ps ...parser) parser {
	return p
}

// zero-or-one matches
func maybe(p parser) parser {
	return p
}

// Name ::= NAME-TYPE TLV-LENGTH NameComponent*
var name = tlv(NameType, many(nameComponent))

// NameComponent ::= GenericNameComponent | ImplicitSha256DigestComponent
var nameComponent = or(genericNameComponent, implicitSha256DigestComponent)

// GenericNameComponent ::= NAME-COMPONENT-TYPE TLV-LENGTH BYTE*
var genericNameComponent = tlv(NameComponentType, bytes)

// ImplicitSha256DigestComponent ::= IMPLICIT-SHA256-DIGEST-COMPONENT-TYPE TLV-LENGTH(=32)
// 																		BYTE{32}
var implicitSha256DigestComponent = tlv(ImplicitSha256DigestComponent, exactBytes(32))
