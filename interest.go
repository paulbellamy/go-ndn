package ndn

type Interest struct {
	name *Name

	childSelector    int
	hasChildSelector bool

	lifetimeMilliseconds    int
	hasLifetimeMilliseconds bool

	maxSuffixComponents    int
	hasMaxSuffixComponents bool

	minSuffixComponents    int
	hasMinSuffixComponents bool

	mustBeFresh    bool
	hasMustBeFresh bool

	scope    int
	hasScope bool
}

func (i *Interest) GetChildSelector() int {
	if !i.hasChildSelector {
		return -1
	}
	return i.childSelector
}

func (i *Interest) SetChildSelector(x int) {
	i.hasChildSelector = true
	i.childSelector = x
}

func (i *Interest) GetInterestLifetimeMilliseconds() int {
	if !i.hasLifetimeMilliseconds {
		return -1
	}
	return i.lifetimeMilliseconds
}

func (i *Interest) SetInterestLifetimeMilliseconds(x int) {
	i.hasLifetimeMilliseconds = true
	i.lifetimeMilliseconds = x
}

func (i *Interest) GetMaxSuffixComponents() int {
	if !i.hasMaxSuffixComponents {
		return -1
	}
	return i.maxSuffixComponents
}

func (i *Interest) SetMaxSuffixComponents(x int) {
	i.hasMaxSuffixComponents = true
	i.maxSuffixComponents = x
}

func (i *Interest) GetMinSuffixComponents() int {
	if !i.hasMinSuffixComponents {
		return -1
	}
	return i.minSuffixComponents
}

func (i *Interest) SetMinSuffixComponents(x int) {
	i.hasMinSuffixComponents = true
	i.minSuffixComponents = x
}

func (i *Interest) GetMustBeFresh() bool {
	if !i.hasMustBeFresh {
		return true
	}
	return i.mustBeFresh
}

func (i *Interest) SetMustBeFresh(x bool) {
	i.hasMustBeFresh = true
	i.mustBeFresh = x
}

func (i *Interest) GetName() *Name {
	if i.name == nil {
		i.name = &Name{}
	}
	return i.name
}

func (i *Interest) SetName(x *Name) {
	i.name = CopyName(x)
}

func (i *Interest) GetScope() int {
	if !i.hasScope {
		return -1
	}
	return i.scope
}

func (i *Interest) SetScope(x int) {
	i.hasScope = true
	i.scope = x
}

func (i *Interest) MatchesName(n *Name) bool {
	return i.GetName().Match(n)
}
