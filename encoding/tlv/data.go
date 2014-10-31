package tlv

import (
	"github.com/paulbellamy/go-ndn/encoding"
	"github.com/paulbellamy/go-ndn/packets"
)

func marshalDataPacket(d *packets.Data) (TLV, error) {
	if d.GetName().IsBlank() {
		return nil, &encoding.NameRequiredError{}
	}

	children := []TLV{marshalName(d.GetName())}

	metaInfoTLV, err := marshalMetaInfo(d.MetaInfo)
	if err != nil {
		return nil, err
	}
	children = append(children, metaInfoTLV)

	contentTLV, err := marshalContent(d.GetContent())
	if err != nil {
		return nil, err
	}
	children = append(children, contentTLV)

	signatureTLVs, err := marshalSignature(d.GetSignature())
	if err != nil {
		return nil, err
	}
	children = append(children, signatureTLVs...)

	return ParentTLV{DataType, children}, nil
}

func marshalMetaInfo(m packets.MetaInfo) (TLV, error) {
	return ParentTLV{MetaInfoType, []TLV{}}, nil
}

func marshalContent(c []byte) (TLV, error) {
	return ByteTLV{ContentType, c}, nil
}

func unmarshalDataPacket(v interface{}, b []byte) error {
	return nil
}
