package ndn

import (
	"github.com/paulbellamy/go-ndn/name"
)

type InterestFilter interface {
	Matches(name.Name) bool
}

type PrefixInterestFilter name.Name

func (f PrefixInterestFilter) Matches(n name.Name) bool {
	return n.GetPrefix(name.Name(f).Size()).Equals(name.Name(f))
}
