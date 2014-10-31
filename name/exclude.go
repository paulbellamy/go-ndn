package name

type Exclude []Component

// A special component used to match anything
var Any = Component{}

// Check if the component matches any of the exclude criteria.
func (e Exclude) Matches(component Component) bool {
	for i := 0; i < len(e); i++ {
		if e[i] == Any {
		} else {
			if component.Equals(e[i]) {
				return true
			}
		}
	}
	return false
}
