package tlv

import "io"

const (
	// Packet Types
	InterestType uint64 = 0x05
	DataType     uint64 = 0x06

	// Common Fields
	NameType                          uint64 = 0x07
	NameComponentType                 uint64 = 0x08
	ImplicitSha256DigestComponentType uint64 = 0x01

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

type valueReader interface {
	io.ReaderFrom
	io.WriterTo
}

type valueReaderFactory func() valueReader

type ChildrenValue []TLV

func (v ChildrenValue) ReadFrom(r io.Reader) (int64, error) {
	return 0, nil
}
func (v ChildrenValue) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
}
func ChildrenReader() valueReader {
	return ChildrenValue{}
}

type UintValue []TLV

func (v UintValue) ReadFrom(r io.Reader) (int64, error) {
	return 0, nil
}
func (v UintValue) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
}
func UintReader() valueReader {
	return UintValue{}
}

type BytesValue []TLV

func (v BytesValue) ReadFrom(r io.Reader) (int64, error) {
	return 0, nil
}
func (v BytesValue) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
}
func BytesReader() valueReader {
	return BytesValue{}
}

var valueReaders = map[uint64]valueReaderFactory{
	// Packet Types
	InterestType: ChildrenReader,
	DataType:     ChildrenReader,

	// Common Fields
	NameType:          ChildrenReader,
	NameComponentType: BytesReader,

	// Interest Packet
	SelectorsType:        ChildrenReader,
	NonceType:            BytesReader,
	ScopeType:            UintReader,
	InterestLifetimeType: UintReader,

	// Interest/Selectors
	MinSuffixComponentsType:       UintReader,
	MaxSuffixComponentsType:       UintReader,
	PublisherPublicKeyLocatorType: ChildrenReader,
	ExcludeType:                   ChildrenReader,
	ChildSelectorType:             UintReader,
	MustBeFreshType:               BytesReader,
	AnyType:                       BytesReader,

	// Data Packet
	MetaInfoType:       ChildrenReader,
	ContentType:        BytesReader,
	SignatureInfoType:  ChildrenReader,
	SignatureValueType: BytesReader,

	// Data/MetaInfo
	ContentTypeType:     UintReader,
	FreshnessPeriodType: UintReader,
	FinalBlockIdType:    ChildrenReader,

	// Data/Signature
	SignatureTypeType: UintReader,
	KeyLocatorType:    ChildrenReader,
	KeyDigestType:     BytesReader,
}
