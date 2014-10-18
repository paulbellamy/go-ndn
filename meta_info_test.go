package ndn

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
