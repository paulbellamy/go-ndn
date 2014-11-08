package packets

import (
	"testing"
	"time"

	"github.com/paulbellamy/go-ndn/name"
	"github.com/stretchr/testify/assert"
)

func Test_MetaInfo_ContentType(t *testing.T) {
	subject := MetaInfo{}

	// when unspecified
	assert.Equal(t, subject.GetContentType(), UnknownContentType)
	subject.SetContentType(LinkContentType)
	assert.Equal(t, subject.GetContentType(), LinkContentType)
}

func Test_MetaInfo_FreshnessPeriod(t *testing.T) {
	subject := MetaInfo{}

	// when unspecified
	assert.Equal(t, subject.GetFreshnessPeriod(), -1)
	subject.SetFreshnessPeriod(2 * time.Minute)
	assert.Equal(t, subject.GetFreshnessPeriod(), 2*time.Minute)

	// check it can be set to zero
	subject.SetFreshnessPeriod(0)
	assert.Equal(t, subject.GetFreshnessPeriod(), 0)
}

func Test_MetaInfo_FinalBlockID(t *testing.T) {
	subject := MetaInfo{}

	// when unspecified
	assert.Equal(t, subject.GetFinalBlockID(), name.Component{})
	subject.SetFinalBlockID(name.Component{"foo"})
	assert.Equal(t, subject.GetFinalBlockID(), name.Component{"foo"})
}
