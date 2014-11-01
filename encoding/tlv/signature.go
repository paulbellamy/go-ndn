package tlv

// Signature ::= SignatureInfo
//               SignatureBits

// SignatureInfo ::= SIGNATURE-INFO-TYPE TLV-LENGTH
//                     SignatureType
//                     KeyLocator?
//                     ... (SignatureType-specific TLVs)

// SignatureValue ::= SIGNATURE-VALUE-TYPE TLV-LENGTH
//                     ... (SignatureType-specific TLVs and
//                     BYTE+
var Signature = or(
	seq(DigestSha256SignatureInfo, DigestSha256SignatureValue),
	seq(SignatureSha256WithRsaInfo, SignatureSha256WithRsaValue),
	seq(SignatureSha256WithEcdsaInfo, SignatureSha256WithEcdsaValue),
)

// SignatureType ::= SIGNATURE-TYPE-TYPE TLV-LENGTH
//                     nonNegativeInteger
var SignatureType = tlv(SignatureTypeType, nonNegativeInteger)

// KeyLocator ::= KEY-LOCATOR-TYPE TLV-LENGTH (Name | KeyDigest)
var KeyLocator = tlv(KeyLocatorType, or(Name, KeyDigest))

// KeyDigest ::= KEY-DIGEST-TYPE TLV-LENGTH BYTE+
var KeyDigest = tlv(KeyDigestType, Bytes)

// DigestSha256

// SignatureInfo ::= SIGNATURE-INFO-TYPE TLV-LENGTH(=3)
//                     SIGNATURE-TYPE-TYPE TLV-LENGTH(=1) 0
var DigestSha256SignatureInfo = tlv(SignatureInfoType,
	tlv(SignatureTypeType, exactly(0)),
)

// SignatureValue ::= SIGNATURE-VALUE-TYPE TLV-LENGTH(=32)
//                      BYTE+(=SHA256{Name, MetaInfo, Content, SignatureInfo})
var DigestSha256SignatureValue = tlv(SignatureValueType,
	Bytes(32),
)

// SignatureSha256WithRsa

// SignatureInfo ::= SIGNATURE-INFO-TYPE TLV-LENGTH
//                     SIGNATURE-TYPE-TYPE TLV-LENGTH(=1) 1
//                     KeyLocator
var SignatureSha256WithRsaInfo = tlv(SignatureInfoType,
	tlv(SignatureTypeType, exactly(1)),
	KeyLocator,
)

// SignatureValue ::= SIGNATURE-VALUE-TYPE TLV-LENGTH
//                      BYTE+(=RSA over SHA256{Name, MetaInfo, Content, SignatureInfo})
var SignatureSha256WithRsaValue = tlv(SignatureValueType,
	Bytes,
)

// SignatureSha256WithEcdsa

// SignatureInfo ::= SIGNATURE-INFO-TYPE TLV-LENGTH
//                     SIGNATURE-TYPE-TYPE TLV-LENGTH(=1) 3
//                     KeyLocator
var SignatureSha256WithEcdsaInfo = tlv(SignatureInfoType,
	tlv(SignatureTypeType, exactly(3)),
	KeyLocator,
)

// SignatureValue ::= SIGNATURE-VALUE-TYPE TLV-LENGTH
//                      BYTE+(=ECDSA over SHA256{Name, MetaInfo, Content, SignatureInfo})
var SignatureSha256WithEcdsaValue = tlv(SignatureValueType,
	Bytes,
)

// Not sure what to do with this
// Ecdsa-Sig-Value  ::=  SEQUENCE  {
//      r     INTEGER,
//      s     INTEGER  }
