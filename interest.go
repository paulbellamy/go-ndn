package ndn

type Interest struct {
	name Name

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

	exclude *Exclude
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
	i.hasMaxSuffixComponents = (x >= 0)
	i.maxSuffixComponents = x
}

func (i *Interest) GetMinSuffixComponents() int {
	if !i.hasMinSuffixComponents {
		return -1
	}
	return i.minSuffixComponents
}

func (i *Interest) SetMinSuffixComponents(x int) {
	i.hasMinSuffixComponents = (x >= 0)
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

func (i *Interest) GetName() Name {
	if i.name == nil {
		i.name = Name{}
	}
	return i.name
}

func (i *Interest) SetName(x Name) {
	i.name = x.Copy()
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

func (i *Interest) GetExclude() *Exclude {
	return i.exclude
}

func (i *Interest) SetExclude(e *Exclude) {
	i.exclude = e
}

// Check if this Interest's name matches the given name and that the given name
// also conforms to the interest selectors.
func (i *Interest) MatchesName(n Name) bool {
	if !i.GetName().Match(n) {
		return false
	}

	if i.hasMinSuffixComponents {
		if i.suffixes(n) < i.minSuffixComponents {
			return false
		}
	}

	if i.hasMaxSuffixComponents {
		if i.suffixes(n) > i.maxSuffixComponents {
			return false
		}
	}

	if i.exclude != nil && n.Size() > i.GetName().Size() {
		if i.exclude.Matches(n.Get(i.GetName().Size())) {
			return false
		}
	}

	return true
}

// Find the number of suffixes in a name
func (i *Interest) suffixes(n Name) int {
	// Add 1 for the implicit digest.
	return n.Size() + 1 - i.GetName().Size()
}
