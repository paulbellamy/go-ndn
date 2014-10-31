package packets

import (
	"testing"
	"time"

	"github.com/paulbellamy/go-ndn/name"
	"github.com/stretchr/testify/assert"
)

func Test_Data_Name(t *testing.T) {
	subject := &Data{}

	// when unspecified
	assert.Equal(t, subject.GetName().Size(), 0)

	n := name.New(name.Component{"a"}, name.Component{"b"})
	subject.SetName(n)
	assert.Equal(t, subject.GetName(), n)

	// check it is copied
	n.Clear()
	assert.Equal(t, subject.GetName(), name.New(name.Component{"a"}, name.Component{"b"}))
}

func Test_Data_FreshnessPeriod(t *testing.T) {
	subject := &Data{}

	// when unspecified
	assert.Equal(t, subject.GetFreshnessPeriod(), -1)
	subject.SetFreshnessPeriod(2 * time.Minute)
	assert.Equal(t, subject.GetFreshnessPeriod(), 2*time.Minute)
}

func Test_Data_Content(t *testing.T) {
	subject := &Data{}

	// when unspecified
	assert.Equal(t, subject.GetContent(), []byte{})
	content := []byte("hello world")
	subject.SetContent(content)
	assert.Equal(t, subject.GetContent(), content)
}

func Test_Data_Signature(t *testing.T) {
	subject := &Data{}

	assert.Nil(t, subject.GetSignature())
	signature := Sha256WithRSASignature{}
	subject.SetSignature(signature)
	assert.Equal(t, subject.GetSignature(), signature)
}
