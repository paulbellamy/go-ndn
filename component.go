package ndn

import (
	"bytes"

	"github.com/paulbellamy/go-ndn/encoding"
)

type Component struct {
	value string
}

func ComponentFromBytes(b []byte) Component {
	return Component{
		value: string(b),
	}
}

func ComponentFromString(s string) Component {
	return Component{
		value: s,
	}
}

func (c Component) Copy() Component {
	return Component{
		value: c.value,
	}
}

func (c Component) String() string {
	return c.value
}

func (c Component) Compare(other Component) int {
	return bytes.Compare([]byte(c.value), []byte(other.value))
}

func (c Component) Equals(other Component) bool {
	return c.value == other.value
}

func (c Component) GetValue() string {
	return c.value
}

func (c Component) toTLV() encoding.TLV {
	return encoding.ByteTLV{encoding.NameComponentType, []byte(c.value)}
}
