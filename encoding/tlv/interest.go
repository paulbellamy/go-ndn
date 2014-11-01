package tlv

// Interest ::= INTEREST-TYPE TLV-LENGTH
//                Name
//                Selectors?
//                Nonce
//                Scope?
//                InterestLifetime?
var Interest = tlv(InterestType,
	Name,
	maybe(Selectors),
	Nonce,
	maybe(Scope),
	maybe(InterestLifetime),
)

// Selectors ::= SELECTORS-TYPE TLV-LENGTH
//                 MinSuffixComponents?
//                 MaxSuffixComponents?
//                 PublisherPublicKeyLocator?
//                 Exclude?
//                 ChildSelector?
//                 MustBeFresh?
var Selectors = tlv(SelectorsType,
	maybe(MinSuffixComponents),
	maybe(MaxSuffixComponents),
	maybe(PublisherPublicKeyLocator),
	maybe(Exclude),
	maybe(ChildSelector),
	maybe(MustBeFresh),
)

// MinSuffixComponents ::= MIN-SUFFIX-COMPONENTS-TYPE TLV-LENGTH
//                           nonNegativeInteger
var MinSuffixComponents = tlv(MinSuffixComponentsType, nonNegativeInteger)

// MaxSuffixComponents ::= MAX-SUFFIX-COMPONENTS-TYPE TLV-LENGTH
//                           nonNegativeInteger
var MaxSuffixComponents = tlv(MaxSuffixComponentsType, nonNegativeInteger)

// PublisherPublicKeyLocator ::= KeyLocator
var PublisherPublicKeyLocator = KeyLocator

// Exclude ::= EXCLUDE-TYPE TLV-LENGTH Any? (NameComponent (Any)?)+
var Exclude = tlv(ExcludeType, maybe(Any), oneOrMore(NameComponent, maybe(Any)))

// Any ::= ANY-TYPE TLV-LENGTH(=0)
var Any = tlv(AnyType, empty)

// ChildSelector ::= CHILD-SELECTOR-TYPE TLV-LENGTH
//                     nonNegativeInteger
var ChildSelector = tlv(ChildSelectorType, nonNegativeInteger)

// MustBeFresh ::= MUST-BE-FRESH-TYPE TLV-LENGTH(=0)
var MustBeFresh = tlv(MustBeFreshType, empty)

// Nonce ::= NONCE-TYPE TLV-LENGTH(=4) BYTE{4}
var Nonce = tlv(NonceType, Bytes(4))

// Guiders

// Scope ::= SCOPE-TYPE TLV-LENGTH nonNegativeInteger
var Scope = tlv(ScopeType, nonNegativeInteger)

// InterestLifetime ::= INTEREST-LIFETIME-TYPE TLV-LENGTH nonNegativeInteger
var InterestLifetime = tlv(InterestLifetimeType, nonNegativeInteger)
