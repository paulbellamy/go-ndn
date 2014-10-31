package ndn

import (
	"testing"

	"github.com/paulbellamy/go-ndn/name"
	"github.com/stretchr/testify/assert"
)

func Test_PrefixInterestFilter(t *testing.T) {
	subject := PrefixInterestFilter(name.New(name.Component{"a"}))
	assert.Implements(t, (*InterestFilter)(nil), subject)

	assert.False(t, subject.Matches(name.New()))
	assert.False(t, subject.Matches(name.New(name.Component{"b"})))
	assert.True(t, subject.Matches(name.New(name.Component{"a"})))
	assert.True(t, subject.Matches(name.New(name.Component{"a"}, name.Component{"b"})))
}
