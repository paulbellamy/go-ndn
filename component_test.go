package ndn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ComponentFromBytes(t *testing.T) {
	assert.Equal(t, &Component{value: "foo"}, ComponentFromBytes([]byte("foo")))
}

func Test_ComponentFromString(t *testing.T) {
	assert.Equal(t, &Component{value: "foo"}, ComponentFromString("foo"))
}

func Test_ComponentFromBlob(t *testing.T) {
	t.Error("pending")
	//assert.Equal(t, &Component{value: "foo"}, ComponentFromBlob("foo"))
}

func Test_ComponentFromNumber(t *testing.T) {
	t.Error("pending")
	//assert.Equal(t, &Component{value: "foo"}, ComponentFromNumber("foo"))
}

func Test_ComponentFromNumberWithMarker(t *testing.T) {
	t.Error("pending")
	//assert.Equal(t, &Component{value: "foo"}, ComponentFromNumberWithMarket("foo"))
}

func Test_CopyComponent(t *testing.T) {
	assert.Equal(t, &Component{value: "foo"}, CopyComponent(ComponentFromString("foo")))
}

func Test_Component_Compare(t *testing.T) {
	subject := &Component{value: "foo"}

	assert.Equal(t, subject.Compare(&Component{value: "fo"}), 1)
	assert.Equal(t, subject.Compare(&Component{value: "fon"}), 1)
	assert.Equal(t, subject.Compare(&Component{value: "foo"}), 0)
	assert.Equal(t, subject.Compare(&Component{value: "fop"}), -1)
	assert.Equal(t, subject.Compare(&Component{value: "fooo"}), -1)
}

func Test_Component_Equals(t *testing.T) {
	subject := &Component{value: "foo"}

	assert.True(t, subject.Equals(&Component{value: "foo"}))
	assert.False(t, subject.Equals(&Component{value: "bar"}))
}

func Test_Component_GetValue(t *testing.T) {
	t.Error("pending")

	//subject := &Component{value: "foo"}
	//assert.Equal(t, subject.GetValue(), &Blob{"foo"})
}

func Test_Component_ToEscapedString(t *testing.T) {
	t.Error("pending")

	//subject := &Component{value: "foo"}
	//assert.Equal(t, subject.GetValue(), &Blob{"foo"})
}

func Test_Component_ToNumber(t *testing.T) {
	t.Error("pending")

	//subject := &Component{value: "123"}
	//assert.Equal(t, subject.ToNumber(), uint64(123))
}

func Test_Component_ToNumberWithMarker(t *testing.T) {
	t.Error("pending")

	//subject := &Component{value: "123"}
	// num, err := subject.ToNumberWithMarket(marker byte)
	//assert.Equal(t, num, uint64(123))
	//assert.NoError(t, err)
	// return error if first byte of the component does not equal the marker
}

func Test_Component_ToSegment(t *testing.T) {
	t.Error("pending")

	//subject := &Component{value: "123"}
	// num, err := subject.ToSegment()
	//assert.Equal(t, num, uint64(123))
	//assert.NoError(t, err)
	// return error if first byte of the component is not the expected marker (0x00)
}

func Test_Component_ToSegmentOffset(t *testing.T) {
	t.Error("pending")

	//subject := &Component{value: "123"}
	// num, err := subject.ToSegmentOffset()
	//assert.Equal(t, num, uint64(123))
	//assert.NoError(t, err)
	// return error if first byte of the component is not the expected marker (0xFB)
}

func Test_Component_ToSequenceNumber(t *testing.T) {
	t.Error("pending")

	//subject := &Component{value: "123"}
	// num, err := subject.ToSequenceNumber()
	//assert.Equal(t, num, uint64(123))
	//assert.NoError(t, err)
	// return error if first byte of the component is not the expected marker (0xFE)
}

func Test_Component_ToTimestamp(t *testing.T) {
	t.Error("pending")

	//subject := &Component{value: "123"}
	// epochUsec, err := subject.ToTimestamp()
	//assert.Equal(t, num, uint64(123))
	//assert.NoError(t, err)
	// return error if first byte of the component is not the expected marker (0xFC)
}

func Test_Component_ToVersion(t *testing.T) {
	t.Error("pending")

	//subject := &Component{value: "123"}
	// epochUsec, err := subject.ToTimestamp()
	//assert.Equal(t, num, uint64(123))
	//assert.NoError(t, err)
	// return error if first byte of the component is not the expected marker (0xFD)
}
