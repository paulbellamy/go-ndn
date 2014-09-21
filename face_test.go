package ndn

import (
	"errors"
	"testing"

	"github.com/paulbellamy/go-ndn/encoding"
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
	pending, err := subject.ExpressInterest(&Interest{name: Name{Component{"a"}}})
	assert.NoError(t, err)
	if assert.NotNil(t, pending) {
		assert.Equal(t, pending.ID, uint64(1))
	}

	assert.Equal(t, len(subject.pendingInterestTable.items), 1)

	transport.AssertExpectations(t)
}

func Test_Face_ExpressInterest_ErrorWriting(t *testing.T) {
	transport := &mockTransport{}
	transport.On("Write", mock.AnythingOfType("[]uint8")).Return(0, errors.New("test error"))

	subject := Face(transport)
	pending, err := subject.ExpressInterest(&Interest{name: Name{Component{"a"}}})
	assert.EqualError(t, err, "test error")
	assert.Nil(t, pending)

	transport.AssertExpectations(t)
}

func Test_Face_RemovePendingInterest(t *testing.T) {
	transport := &mockTransport{}
	transport.On("Write", mock.AnythingOfType("[]uint8")).Return(0, nil)
	subject := Face(transport)

	// Add one into the table
	pending, err := subject.ExpressInterest(&Interest{name: Name{Component{"a"}}})
	assert.NoError(t, err)
	assert.Equal(t, len(subject.pendingInterestTable.items), 1)

	// Remove it and check it's gone
	subject.RemovePendingInterest(pending.ID)
	assert.Equal(t, len(subject.pendingInterestTable.items), 0)

	transport.AssertExpectations(t)
}

func Test_Face_ReceivingData(t *testing.T) {
	transport := &bufferTransport{}
	subject := Face(transport)

	// Add one into the table
	pending, err := subject.ExpressInterest(&Interest{name: Name{Component{"a"}}})
	assert.NoError(t, err)
	assert.Equal(t, len(subject.pendingInterestTable.items), 1)

	// Make some data available
	_, err = encoding.ParentTLV{
		T: encoding.DataType,
		V: []encoding.TLV{
			Name{Component{"a"}}.toTLV(),
		},
	}.WriteTo(transport)
	assert.NoError(t, err)
	transport.Close()

	// Process the data, EOF is silenced
	assert.NoError(t, subject.ProcessEvents())

	// Check we received it
	select {
	case d, ok := <-pending.Data:
		assert.True(t, ok)
		assert.NotNil(t, d)
	default:
		t.Error("Timeout waiting for data packet")
	}

	// Check the channel was closed
	select {
	case d, ok := <-pending.Data:
		assert.False(t, ok)
		assert.Nil(t, d)
	default:
		t.Error("Timeout waiting for channel to close")
	}
}
