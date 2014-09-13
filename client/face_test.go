package client

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Face(t *testing.T) {
	transport := testTransport(t)
	subject := Face(transport)
	assert.IsType(t, (*face)(nil), subject)
}

func Test_Face_ExpressInterest(t *testing.T) {
	transport := &mockTransport{}
	transport.On("Write", mock.AnythingOfType("[]uint8")).Return(0, nil)

	subject := Face(transport)
	pending, err := subject.ExpressInterest(&Interest{Name: Name{"a"}})
	assert.NoError(t, err)
	if assert.NotNil(t, pending) {
		// TODO: check stuff about the pending here
	}

	assert.Equal(t, len(subject.pendingInterestTable.items), 1)

	transport.AssertExpectations(t)
}

func Test_Face_ExpressInterest_ErrorWriting(t *testing.T) {
	transport := &mockTransport{}
	transport.On("Write", mock.AnythingOfType("[]uint8")).Return(0, errors.New("test error"))

	subject := Face(transport)
	pending, err := subject.ExpressInterest(&Interest{Name: Name{"a"}})
	assert.EqualError(t, err, "test error")
	assert.Nil(t, pending)

	transport.AssertExpectations(t)
}
