package name

import (
	"errors"
	"fmt"
	"strings"
)

var ErrTLVIsNotAName = errors.New("TLV is not a name")

type Name []Component

func New(cs ...Component) Name {
	return Name(cs)
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

func (n Name) ToURI() string {
	segmentStrings := []string{}
	for _, segment := range n {
		segmentStrings = append(segmentStrings, segment.ToURI())
	}
	return fmt.Sprint("ndn:/", strings.Join(segmentStrings, "/"))
}
