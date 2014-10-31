package tlv

import (
	"github.com/paulbellamy/go-ndn/packets"
)

func marshalSignature(s packets.Signature) ([]TLV, error) {
	if s == nil {
		return []TLV{}, nil
	}
	return []TLV{
		ParentTLV{
			SignatureInfoType,
			[]TLV{
				UintTLV{SignatureTypeType, s.Type()},
				// optional KeyLocator
			},
		},
		ByteTLV{SignatureValueType, s.Bytes()},
	}, nil
}
