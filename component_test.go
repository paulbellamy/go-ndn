package ndn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ComponentFromBytes(t *testing.T) {
	assert.Equal(t, Component{value: "foo"}, ComponentFromBytes([]byte("foo")))
}

func Test_ComponentFromString(t *testing.T) {
	assert.Equal(t, Component{value: "foo"}, ComponentFromString("foo"))
}

func Test_Component_Copy(t *testing.T) {
	assert.Equal(t, Component{value: "foo"}, ComponentFromString("foo").Copy())
}

func Test_Component_Compare(t *testing.T) {
	subject := Component{value: "foo"}

	assert.Equal(t, subject.Compare(Component{value: "fo"}), 1)
	assert.Equal(t, subject.Compare(Component{value: "fon"}), 1)
	assert.Equal(t, subject.Compare(Component{value: "foo"}), 0)
	assert.Equal(t, subject.Compare(Component{value: "fop"}), -1)
	assert.Equal(t, subject.Compare(Component{value: "fooo"}), -1)
}

func Test_Component_Equals(t *testing.T) {
	subject := Component{value: "foo"}

	assert.True(t, subject.Equals(Component{value: "foo"}))
	assert.False(t, subject.Equals(Component{value: "bar"}))
}

func Test_Component_GetValue(t *testing.T) {
	assert.Equal(t, Component{value: "foo"}.GetValue(), "foo")
}

func Test_Component_ToURI_Escaping(t *testing.T) {
	assert.Equal(t, Component{"b c"}.ToURI(), "b%20c")
}

func Test_Component_ToURI_Empty(t *testing.T) {
	assert.Equal(t, Component{""}.ToURI(), "...")
}

func Test_Component_ToURI_Dots(t *testing.T) {
	assert.Equal(t, Component{".."}.ToURI(), ".....")
}
