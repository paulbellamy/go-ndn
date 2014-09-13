package encoding

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Interest_Read_Minimal(t *testing.T) {
	subject := &Interest{}
	buf := &bytes.Buffer{}
	assert.NoError(t, (&TLV{Type: NameType, Value: []byte("name1")}).WriteTo(buf))
	assert.NoError(t, (&UintTLV{Type: NonceType, Value: 1234}).WriteTo(buf))
	n, err := subject.ReadFrom(buf)
	assert.NoError(t, err)
	assert.Equal(t, n, 7+6)
	assert.Equal(t, subject.Name, "name1")
	assert.Equal(t, subject.Nonce, uint64(1234))
}

/*
func Test_Interest_Read_Maximal(t *testing.T) {
	subject := &Interest{}
	subject.ReadFrom(bytes.NewReader([]byte{
	// Name
	// Selectors
	// Nonce
	// Scope
	// InterestLifetime
	}))
}
*/

func Test_Interest_Read_MalformedPacket(t *testing.T) {
	subject := &Interest{}
	buf := &bytes.Buffer{}
	assert.NoError(t, (&TLV{Type: SelectorsType, Value: []byte("name1")}).WriteTo(buf))
	assert.NoError(t, (&UintTLV{Type: AnyType, Value: 1234}).WriteTo(buf))

	n, err := subject.ReadFrom(buf)
	assert.Equal(t, n, 7)
	assert.EqualError(t, err, ErrUnexpectedTLVType.Error())
}

func Test_Interest_Read_ErrorHandling(t *testing.T) {
	r := &mockReader{}
	r.On("Read", mock.AnythingOfType("[]uint8")).Return(12, errors.New("test error"))

	subject := &Interest{}
	n, err := subject.ReadFrom(r)
	assert.EqualError(t, err, "test error")
	assert.Equal(t, n, 12)

	r.AssertExpectations(t)
}

type mockReader struct {
	mock.Mock
}

func (m *mockReader) Read(b []byte) (int, error) {
	args := m.Mock.Called(b)
	return args.Int(0), args.Error(1)
}
