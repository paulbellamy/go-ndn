package ndn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseURI(t *testing.T) {
	name, err := ParseURI("ndn:/a/b")
	assert.NoError(t, err)
	assert.Equal(t, name, Name{Component{"a"}, Component{"b"}})
}

func Test_ParseURI_WrongProtocol(t *testing.T) {
	name, err := ParseURI("http://a+._-c/b%20c")
	assert.EqualError(t, err, "protocol scheme must be ndn")
	assert.Nil(t, name)
}

func Test_ParseURI_Unescaping(t *testing.T) {
	name, err := ParseURI("ndn:/a+._-c/b%20c")
	assert.NoError(t, err)
	assert.Equal(t, name, Name{Component{"a+._-c"}, Component{"b c"}})
}

func Test_ParseURI_EqualsUnescaping(t *testing.T) {
	t.Error("pending")
}

func Test_ParseURI_LeadingDoubleSlash(t *testing.T) {
	name, err := ParseURI("ndn://hostname/a/b")
	assert.NoError(t, err)
	assert.Equal(t, name, Name{Component{"a"}, Component{"b"}})
}

func Test_ParseURI_NoLeadingSlash(t *testing.T) {
	name, err := ParseURI("ndn:a/b")
	assert.EqualError(t, err, "invalid uri")
	assert.Nil(t, name)
}

func Test_ParseURI_LeadingTripleSlash(t *testing.T) {
	// Not sure what we are supposed to do here for NDN. Just treating at as an empty hostname for now.
	name, err := ParseURI("ndn:///a/b")
	assert.NoError(t, err)
	assert.Equal(t, name, Name{Component{"a"}, Component{"b"}})
}

func Test_ParseURI_IncorrectEmptySegments(t *testing.T) {
	name, err := ParseURI("ndn:/a//b")
	assert.EqualError(t, err, "invalid URI escape \"\"")
	assert.Nil(t, name)
}

func Test_ParseURI_EmptyEscapedSegments(t *testing.T) {
	name, err := ParseURI("ndn:/a/.../b")
	assert.NoError(t, err)
	assert.Equal(t, name, Name{Component{"a"}, Component{""}, Component{"b"}})
}

func Test_ParseURI_DottedSegments(t *testing.T) {
	name, err := ParseURI("ndn:/a/...../b")
	assert.NoError(t, err)
	assert.Equal(t, name, Name{Component{"a"}, Component{".."}, Component{"b"}})
}

func Test_ParseURI_DottedSegments_IncorrectEscaping(t *testing.T) {
	name, err := ParseURI("ndn:/a/../b")
	assert.EqualError(t, err, "invalid URI escape \"..\"")
	assert.Nil(t, name)
}
