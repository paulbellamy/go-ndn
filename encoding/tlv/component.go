package tlv

import (
	"github.com/paulbellamy/go-ndn/name"
)

func marshalComponent(c name.Component) TLV {
	return ByteTLV{NameComponentType, c.Bytes()}
}
