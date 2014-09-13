package ndn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Interest_ChildSelector(t *testing.T) {
	subject := &Interest{}

	// when unspecified
	assert.Equal(t, subject.GetChildSelector(), -1)

	subject.SetChildSelector(2)
	assert.Equal(t, subject.GetChildSelector(), 2)
}

func Test_Interest_InterestLifetimeMilliseconds(t *testing.T) {
	subject := &Interest{}

	// when unspecified
	assert.Equal(t, subject.GetInterestLifetimeMilliseconds(), -1)

	subject.SetInterestLifetimeMilliseconds(2)
	assert.Equal(t, subject.GetInterestLifetimeMilliseconds(), 2)
}

func Test_Interest_MaxSuffixComponents(t *testing.T) {
	subject := &Interest{}

	// when unspecified
	assert.Equal(t, subject.GetMaxSuffixComponents(), -1)

	subject.SetMaxSuffixComponents(2)
	assert.Equal(t, subject.GetMaxSuffixComponents(), 2)
}

func Test_Interest_MinSuffixComponents(t *testing.T) {
	subject := &Interest{}

	// when unspecified
	assert.Equal(t, subject.GetMinSuffixComponents(), -1)

	subject.SetMinSuffixComponents(2)
	assert.Equal(t, subject.GetMinSuffixComponents(), 2)
}

func Test_Interest_MustBeFresh(t *testing.T) {
	subject := &Interest{}
	// when unspecified
	assert.True(t, subject.GetMustBeFresh())

	subject.SetMustBeFresh(true)
	assert.True(t, subject.GetMustBeFresh())

	subject.SetMustBeFresh(false)
	assert.False(t, subject.GetMustBeFresh())
}

func Test_Interest_Name(t *testing.T) {
	subject := &Interest{}

	// when unspecified
	assert.Equal(t, subject.GetName().Size(), 0)

	name := Name{Component{"a"}, Component{"b"}}
	subject.SetName(name)
	assert.Equal(t, subject.GetName(), name)

	// check it is copied
	name.Clear()
	assert.Equal(t, subject.GetName(), Name{Component{"a"}, Component{"b"}})
}

func Test_Interest_Scope(t *testing.T) {
	subject := &Interest{}

	// when unspecified
	assert.Equal(t, subject.GetScope(), -1)

	subject.SetScope(2)
	assert.Equal(t, subject.GetScope(), 2)
}

func Test_Interest_Exclude(t *testing.T) {
	subject := &Interest{}

	// when unspecified
	assert.Nil(t, subject.GetExclude())

	e := &Exclude{}
	subject.SetExclude(e)
	assert.Equal(t, subject.GetExclude(), e)
}

func Test_Interest_MatchesName(t *testing.T) {
	subject := &Interest{}
	subject.SetName(Name{Component{"foo"}})
	assert.True(t, subject.MatchesName(Name{Component{"foo"}, Component{"bar"}}))

	subject.SetMinSuffixComponents(1)
	assert.False(t, subject.MatchesName(Name{}))
	subject.SetMinSuffixComponents(-1)

	subject.SetMaxSuffixComponents(1)
	assert.False(t, subject.MatchesName(Name{Component{"foo"}, Component{"bar"}, Component{"baz"}}))
	subject.SetMaxSuffixComponents(-1)

	subject.SetExclude(&Exclude{Component{"bar"}})
	assert.False(t, subject.MatchesName(Name{Component{"foo"}, Component{"bar"}}))
	subject.SetExclude(nil)
}
