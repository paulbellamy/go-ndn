package ndn

type InterestFilter interface {
	Matches(Name) bool
}

type PrefixInterestFilter Name

func (f PrefixInterestFilter) Matches(n Name) bool {
	return n.GetPrefix(Name(f).Size()).Equals(Name(f))
}
