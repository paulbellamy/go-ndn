package ndn

import (
	"fmt"
	"strings"
)

type Name []*Component

func CopyName(n *Name) *Name {
	newName := &Name{}
	for _, component := range *n {
		newName.AppendComponent(component)
	}
	return newName
}

// implement this with Write as well
func (n *Name) AppendBytes(b []byte) *Name {
	return n
}

func (n *Name) AppendString(s string) *Name {
	return n
}

func (n *Name) AppendComponent(c *Component) *Name {
	(*n) = append(*n, c)
	return n
}

func (n *Name) Clear() {
	(*n) = []*Component{}
}

func (n *Name) Get(i int) *Component {
	size := n.Size()

	if i >= size || i < size*-1 {
		return nil
	} else if i < 0 {
		return (*n)[size+i]
	}
	return (*n)[i]
}

func (n *Name) Set(uri string) {
}

func (n *Name) GetPrefix(i int) *Name {
	newName := &Name{}
	size := n.Size()
	if i < 0 {
		i += size
	}
	for x := 0; x < i && x < size; x++ {
		newName.AppendComponent(n.Get(x))
	}
	return newName
}

func (n *Name) GetSubName(offset, count int) *Name {
	newName := &Name{}
	size := n.Size()
	if offset < 0 {
		offset += size
	}
	for x := offset; x < offset+count && x < size; x++ {
		newName.AppendComponent(n.Get(x))
	}
	return newName
}

func (n *Name) Compare(other *Name) int {
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

func (n *Name) Match(other *Name) bool {
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

func (n *Name) Size() int {
	return len(*n)
}

func (n *Name) ToURI() string {
	return ""
}

func (n *Name) String() string {
	segments := []string{}
	for _, component := range *n {
		segments = append(segments, component.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(segments, ", "))
}
