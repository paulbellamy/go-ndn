package tlv

import (
	"testing"

	"github.com/paulbellamy/go-ndn/name"
	"github.com/stretchr/testify/assert"
)

func Test_marshalComponent(t *testing.T) {
	assert.Equal(t, marshalComponent(name.Component{"foo"}), ByteTLV{
		T: NameComponentType,
		V: []byte("foo"),
	})
}
