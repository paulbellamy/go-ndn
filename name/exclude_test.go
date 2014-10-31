package name

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Exclude_MatchesComponents(t *testing.T) {
	assert.False(t, Exclude{Component{"a"}}.Matches(Component{"b"}))
	assert.True(t, Exclude{Component{"a"}}.Matches(Component{"a"}))
}

func Test_Exclude_AnyCriteria(t *testing.T) {
	t.Error("I don't understand this")
}
