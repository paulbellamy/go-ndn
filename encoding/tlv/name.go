package tlv

import (
	"io"

	"github.com/paulbellamy/go-ndn/encoding"
	"github.com/paulbellamy/go-ndn/name"
)

func marshalName(n name.Name) TLV {
	componentTLVs := []TLV{}
	for _, component := range n {
		componentTLVs = append(componentTLVs, marshalComponent(component))
	}
	return ParentTLV{NameType, componentTLVs}
}

func unmarshalName(r io.Reader) (n name.Name, err error) {
	parentTLV, err := readParentTLV(r)
	if err != nil {
		return n, err
	}

	if parentTLV.T != NameType {
		return n, &encoding.InvalidUnmarshalError{Message: "name tlv expected"}
	}

	for _, component := range parentTLV.V {
		component, ok := component.(ByteTLV)
		if !ok || component.T != NameComponentType {
			return n, &encoding.InvalidUnmarshalError{Message: "name component tlv expected"}
		}

		n = n.AppendBytes(component.V)
	}
	return n, err
}

/*
func NameFromTLV(t tlv.TLV) (n Name, err error) {
	if t.Type() != tlv.NameType {
		err = ErrTLVIsNotAName
		return
	}
	parent, ok := t.(tlv.ParentTLV)
	if !ok {
		err = ErrTLVIsNotAName
		return
	}

	for _, component := range parent.V {
		c, ok := component.(tlv.ByteTLV)
		if !ok {
			err = ErrTLVIsNotAName
			return
		}

		n = append(n, Component{string(c.V)})
	}

	return
}

func (n Name) toTLV() tlv.TLV {
	componentTLVs := []tlv.TLV{}
	for _, component := range n {
		componentTLVs = append(componentTLVs, component.toTLV())
	}
	return tlv.ParentTLV{tlv.NameType, componentTLVs}
}
*/
