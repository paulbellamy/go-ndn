package name

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ComponentFromBytes(t *testing.T) {
	assert.Equal(t, Component{Value: "foo"}, ComponentFromBytes([]byte("foo")))
}

func Test_ComponentFromString(t *testing.T) {
	assert.Equal(t, Component{Value: "foo"}, ComponentFromString("foo"))
}

func Test_Component_Copy(t *testing.T) {
	assert.Equal(t, Component{Value: "foo"}, ComponentFromString("foo").Copy())
}

func Test_Component_Compare(t *testing.T) {
	subject := Component{Value: "foo"}

	assert.Equal(t, subject.Compare(Component{Value: "fo"}), 1)
	assert.Equal(t, subject.Compare(Component{Value: "fon"}), 1)
	assert.Equal(t, subject.Compare(Component{Value: "foo"}), 0)
	assert.Equal(t, subject.Compare(Component{Value: "fop"}), -1)
	assert.Equal(t, subject.Compare(Component{Value: "fooo"}), -1)
}

func Test_Component_Equals(t *testing.T) {
	subject := Component{Value: "foo"}

	assert.True(t, subject.Equals(Component{Value: "foo"}))
	assert.False(t, subject.Equals(Component{Value: "bar"}))
}

func Test_Component_GetValue(t *testing.T) {
	assert.Equal(t, Component{Value: "foo"}.GetValue(), "foo")
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
