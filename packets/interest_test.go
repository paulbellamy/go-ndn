package packets

import (
	"testing"
	"time"

	"github.com/paulbellamy/go-ndn/name"
	"github.com/stretchr/testify/assert"
)

func Test_Interest_ChildSelector(t *testing.T) {
	subject := &Interest{}

	// when unspecified
	assert.Equal(t, subject.GetChildSelector(), -1)

	subject.SetChildSelector(2)
	assert.Equal(t, subject.GetChildSelector(), 2)
}

func Test_Interest_InterestLifetime(t *testing.T) {
	subject := &Interest{}

	// when unspecified
	assert.Equal(t, subject.GetInterestLifetime(), -1)

	subject.SetInterestLifetime(2 * time.Millisecond)
	assert.Equal(t, subject.GetInterestLifetime(), 2*time.Millisecond)
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

func Test_Interest_Signature(t *testing.T) {
	subject := &Interest{}

	assert.Nil(t, subject.GetSignature())
	signature := Sha256WithRSASignature{}
	subject.SetSignature(signature)
	assert.Equal(t, subject.GetSignature(), signature)
}

func Test_Interest_Name(t *testing.T) {
	subject := &Interest{}

	// when unspecified
	assert.Equal(t, subject.GetName().Size(), 0)

	n := name.New(name.Component{"a"}, name.Component{"b"})
	subject.SetName(n)
	assert.Equal(t, subject.GetName(), n)

	// check it is copied
	n.Clear()
	assert.Equal(t, subject.GetName(), name.New(name.Component{"a"}, name.Component{"b"}))
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

	e := &name.Exclude{}
	subject.SetExclude(e)
	assert.Equal(t, subject.GetExclude(), e)
}

func Test_Interest_MatchesName(t *testing.T) {
	subject := &Interest{}
	subject.SetName(name.New(name.Component{"foo"}))
	assert.True(t, subject.MatchesName(name.New(name.Component{"foo"}, name.Component{"bar"})))

	subject.SetMinSuffixComponents(1)
	assert.False(t, subject.MatchesName(name.New()))
	subject.SetMinSuffixComponents(-1)

	subject.SetMaxSuffixComponents(1)
	assert.False(t, subject.MatchesName(name.New(name.Component{"foo"}, name.Component{"bar"}, name.Component{"baz"})))
	subject.SetMaxSuffixComponents(-1)

	subject.SetExclude(&name.Exclude{name.Component{"bar"}})
	assert.False(t, subject.MatchesName(name.New(name.Component{"foo"}, name.Component{"bar"})))
	subject.SetExclude(nil)
}
