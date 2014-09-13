package client

import (
	"github.com/paulbellamy/go-ndn/encoding"
)

type Name []string

func (n Name) IsBlank() bool {
	return len(n) == 0
}

func (n Name) toTLV() encoding.TLV {
	componentTLVs := []encoding.TLV{}
	for _, component := range n {
		componentTLVs = append(componentTLVs, encoding.ByteTLV(encoding.NameComponentType, []byte(component)))
	}
	return encoding.ParentTLV(encoding.NameType, componentTLVs...)
}
