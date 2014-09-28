package ndn

import (
	"errors"
	"fmt"
	"strings"

	"github.com/paulbellamy/go-ndn/encoding"
)

var ErrTLVIsNotAName = errors.New("TLV is not a name")

type Name []Component

func NameFromTLV(t encoding.TLV) (n Name, err error) {
	if t.Type() != encoding.NameType {
		err = ErrTLVIsNotAName
		return
	}
	tlv, ok := t.(encoding.ParentTLV)
	if !ok {
		err = ErrTLVIsNotAName
		return
	}

	for _, component := range tlv.V {
		c, ok := component.(encoding.ByteTLV)
		if !ok {
			err = ErrTLVIsNotAName
			return
		}

		n = append(n, Component{string(c.V)})
	}

	return
}

func (n Name) Copy() Name {
	newName := Name{}
	for _, component := range n {
		newName = append(newName, component)
	}
	return newName
}

// implement this with Write as well
func (n Name) AppendBytes(b []byte) Name {
	return append(n, ComponentFromBytes(b))
}

func (n Name) AppendString(s string) Name {
	return n
}

func (n Name) AppendComponent(c Component) Name {
	return append(n, c)
}

func (n *Name) Clear() {
	(*n) = []Component{}
}

func (n Name) Get(i int) Component {
	return n[n.normalizeIndex(i)]
}

func (n Name) GetPrefix(count int) Name {
	return n.GetSubName(0, n.normalizeIndex(count))
}

func (n Name) GetSuffix(count int) Name {
	start := n.Size() - count
	if start < 0 {
		start = 0
	}
	return n.GetSubName(start, n.Size())
}

func (n Name) GetSubName(offset, count int) Name {
	start := n.normalizeIndex(offset)
	end := start + count
	if end > n.Size() {
		end = n.Size()
	}
	return n[start:end]
}

func (n Name) normalizeIndex(i int) int {
	size := n.Size()
	if i < 0 {
		return size + i
	}

	return i
}

func (n Name) Compare(other Name) int {
	nSize := n.Size()
	otherSize := other.Size()
	for i := 0; i < nSize && i < otherSize; i++ {
		result := n.Get(i).Compare(other.Get(i))
		if result != 0 {
			return result
		}
	}

	if nSize < otherSize {
		return -1
	} else if otherSize < nSize {
		return 1
	}

	return 0
}

func (n Name) Equals(other Name) bool {
	return n.Compare(other) == 0
}

func (n Name) Match(other Name) bool {
	toCheck := n.Size()
	if other.Size() < toCheck {
		return false
	}

	for x := 0; x < toCheck; x++ {
		if n.Get(x).Compare(other.Get(x)) != 0 {
			return false
		}
	}

	return true
}

func (n Name) Size() int {
	return len(n)
}

func (n Name) IsEmpty() bool {
	return n.Size() == 0
}

func (n Name) String() string {
	segments := []string{}
	for _, component := range n {
		segments = append(segments, component.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(segments, ", "))
}

func (n Name) IsBlank() bool {
	return len(n) == 0
}

func (n Name) toTLV() encoding.TLV {
	componentTLVs := []encoding.TLV{}
	for _, component := range n {
		componentTLVs = append(componentTLVs, component.toTLV())
	}
	return encoding.ParentTLV{encoding.NameType, componentTLVs}
}

func (n Name) ToURI() string {
	segmentStrings := []string{}
	for _, segment := range n {
		segmentStrings = append(segmentStrings, escapeDots(escapeHex(segment.String())))
	}
	return fmt.Sprint("ndn:/", strings.Join(segmentStrings, "/"))
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
