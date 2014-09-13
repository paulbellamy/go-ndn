package ndn

import "bytes"

type Component struct {
	value string
}

func ComponentFromBytes(b []byte) *Component {
	return &Component{
		value: string(b),
	}
}

func ComponentFromString(s string) *Component {
	return &Component{
		value: s,
	}
}

func CopyComponent(c *Component) *Component {
	return &Component{
		value: c.value,
	}
}

func (c *Component) String() string {
	return c.value
}

func (c *Component) Compare(other *Component) int {
	return bytes.Compare([]byte(c.value), []byte(other.value))
}

func (c *Component) Equals(other *Component) bool {
	return c.value == other.value
}
