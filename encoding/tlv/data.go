package tlv

// Data ::= DATA-TLV TLV-LENGTH
//            Name
//            MetaInfo
//            Content
//            Signature
var Data = tlv(DataType,
	Name,
	MetaInfo,
	Content,
	Signature,
)

// MetaInfo ::= META-INFO-TYPE TLV-LENGTH
//                ContentType?
//                FreshnessPeriod?
//                FinalBlockId?
var MetaInfo = tlv(MetaInfoType,
	maybe(ContentTypeParser),
	maybe(FreshnessPeriod),
	maybe(FinalBlockID),
)

// ContentType ::= CONTENT-TYPE-TYPE TLV-LENGTH
//                   nonNegativeInteger
var ContentTypeParser = tlv(ContentTypeType, nonNegativeInteger)

// FreshnessPeriod ::= FRESHNESS-PERIOD-TLV TLV-LENGTH
//                       nonNegativeInteger
var FreshnessPeriod = tlv(FreshnessPeriodType, nonNegativeInteger)

// FinalBlockId ::= FINAL-BLOCK-ID-TLV TLV-LENGTH
//                       NameComponent
var FinalBlockID = tlv(FinalBlockIdType, NameComponent)

// Content ::= CONTENT-TYPE TLV-LENGTH BYTE*
var Content = tlv(ContentType, Bytes)
