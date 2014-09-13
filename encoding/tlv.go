package encoding

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

type TLV struct {
	Type  uint64
	Value []byte
}
