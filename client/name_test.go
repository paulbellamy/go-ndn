package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Name_IsBlank(t *testing.T) {
	assert.True(t, Name{}.IsBlank())
	assert.False(t, Name{"a"}.IsBlank())
}
