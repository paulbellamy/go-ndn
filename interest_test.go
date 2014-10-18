package ndn

import (
	"bytes"
	"testing"
	"time"

	"github.com/paulbellamy/go-ndn/encoding"
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

func Test_Interest_WriteTo(t *testing.T) {
	expected, err := encoding.ParentTLV{
		T: encoding.InterestType,
		V: []encoding.TLV{
			encoding.ParentTLV{
				T: encoding.NameType,
				V: []encoding.TLV{
					encoding.ByteTLV{
						T: encoding.NameComponentType,
						V: []byte{'a'},
					},
				},
			},
		},
	}.MarshalBinary()
	assert.NoError(t, err)

	buf := &bytes.Buffer{}
	subject := &Interest{
		name: Name{Component{"a"}},
	}
	n, err := subject.WriteTo(buf)
	assert.NoError(t, err)
	assert.Equal(t, n, int64(len(buf.Bytes())))
	assert.Equal(t, buf.Bytes(), expected)
}

func Test_Interest_WriteTo_WithoutAName(t *testing.T) {
	buf := &bytes.Buffer{}

	subject := &Interest{}
	n, err := subject.WriteTo(buf)
	assert.EqualError(t, err, ErrNameRequired.Error())
	assert.Equal(t, n, int64(0))
	assert.Nil(t, buf.Bytes(), "Expected nothing to be written, but data was found.")
}
