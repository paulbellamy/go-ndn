package ndn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PrefixInterestFilter(t *testing.T) {
	subject := PrefixInterestFilter(Name{Component{"a"}})
	assert.Implements(t, (*InterestFilter)(nil), subject)

	assert.False(t, subject.Matches(Name{}))
	assert.False(t, subject.Matches(Name{Component{"b"}}))
	assert.True(t, subject.Matches(Name{Component{"a"}}))
	assert.True(t, subject.Matches(Name{Component{"a"}, Component{"b"}}))
}
