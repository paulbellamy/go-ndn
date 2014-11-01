package tlv

// Name ::= NAME-TYPE TLV-LENGTH NameComponent*
var Name = tlv(NameType, zeroOrMore(NameComponent))

// NameComponent ::= GenericNameComponent | ImplicitSha256DigestComponent
var NameComponent = or(GenericNameComponent, ImplicitSha256DigestComponent)

// GenericNameComponent ::= NAME-COMPONENT-TYPE TLV-LENGTH BYTE*
var GenericNameComponent = tlv(NameComponentType, Bytes)

// ImplicitSha256DigestComponent ::= IMPLICIT-SHA256-DIGEST-COMPONENT-TYPE TLV-LENGTH(=32)
// 																		BYTE{32}
var ImplicitSha256DigestComponent = tlv(ImplicitSha256DigestComponentType, Bytes(32))
