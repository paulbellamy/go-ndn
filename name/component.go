package name

import (
	"bytes"
	"fmt"
)

type Component struct {
	Value string
}

func ComponentFromBytes(b []byte) Component {
	return Component{
		Value: string(b),
	}
}

func ComponentFromString(s string) Component {
	return Component{
		Value: s,
	}
}

func (c Component) Copy() Component {
	return Component{
		Value: c.Value,
	}
}

func (c Component) String() string {
	return c.Value
}

func (c Component) Bytes() []byte {
	return []byte(c.Value)
}

func (c Component) Compare(other Component) int {
	return bytes.Compare([]byte(c.Value), []byte(other.Value))
}

func (c Component) Equals(other Component) bool {
	return c.Value == other.Value
}

func (c Component) GetValue() string {
	return c.Value
}

func (c Component) ToURI() string {
	return escapeDots(escapeHex(c.String()))
}

func escapeHex(s string) string {
	spaceCount, hexCount := 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c) {
			hexCount++
		}
	}

	if spaceCount == 0 && hexCount == 0 {
		return s
	}

	t := make([]byte, len(s)+2*hexCount)
	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case shouldEscape(c):
			t[j] = '%'
			t[j+1] = "0123456789ABCDEF"[c>>4]
			t[j+2] = "0123456789ABCDEF"[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}

// Return true if the specified character should be escaped when
// appearing in a URL string, according to RFC 3986.
// When 'all' is true the full range of reserved characters are matched.
func shouldEscape(c byte) bool {
	// §2.3 Unreserved characters (alphanum)
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
		return false
	}

	switch c {
	case '-', '_', '.', '~': // §2.3 Unreserved characters (mark)
		return false

	case '$', '&', '+', ',', '/', ':', ';', '=', '?', '@': // §2.2 Reserved characters (reserved)
		// The RFC allows : @ & = + $ but saves / ; , for assigning
		// meaning to individual path segments. This package
		// only manipulates the path as a whole, so we allow those
		// last two as well. That leaves only ? to escape.
		return c == '?'
	}

	// Everything else must be escaped.
	return true
}

func escapeDots(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] != '.' {
			return s
		}
	}

	return fmt.Sprint("...", s)
}
