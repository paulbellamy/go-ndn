package ndn

import (
	"errors"
	"testing"

	"github.com/paulbellamy/go-ndn/encoding/tlv"
	"github.com/paulbellamy/go-ndn/name"
	"github.com/paulbellamy/go-ndn/packets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Face(t *testing.T) {
	transport := testTransport(t)
	subject := NewFace(transport)
	assert.IsType(t, (*Face)(nil), subject)
}

func Test_Face_ExpressInterest(t *testing.T) {
	transport := &mockTransport{}
	transport.On("Write", mock.AnythingOfType("[]uint8")).Return(0, nil)

	packet := &packets.Interest{}
	packet.SetName(name.New(name.Component{"a"}))

	subject := NewFace(transport)
	pending, err := subject.ExpressInterest(packet)
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

	packet := &packets.Interest{}
	packet.SetName(name.New(name.Component{"a"}))

	subject := NewFace(transport)
	pending, err := subject.ExpressInterest(packet)
	assert.EqualError(t, err, "test error")
	assert.Nil(t, pending)

	transport.AssertExpectations(t)
}

func Test_Face_RemovePendingInterest(t *testing.T) {
	transport := &mockTransport{}
	transport.On("Write", mock.AnythingOfType("[]uint8")).Return(0, nil)
	subject := NewFace(transport)

	packet := &packets.Interest{}
	packet.SetName(name.New(name.Component{"a"}))

	// Add one into the table
	pending, err := subject.ExpressInterest(packet)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, len(subject.pendingInterestTable.items), 1)

	// Remove it and check it's gone
	subject.RemovePendingInterest(pending.ID)
	assert.Equal(t, len(subject.pendingInterestTable.items), 0)

	transport.AssertExpectations(t)
}

func Test_Face_ReceivingData(t *testing.T) {
	transport := &bufferTransport{}
	subject := NewFace(transport)

	packet := &packets.Interest{}
	packet.SetName(name.New(name.Component{"a"}))

	// Add one into the table
	pending, err := subject.ExpressInterest(packet)
	assert.NoError(t, err)
	assert.Equal(t, len(subject.pendingInterestTable.items), 1)

	// Make some data available
	data := &packets.Data{}
	data.SetName(name.New(name.Component{"a"}))
	assert.NoError(t, tlv.NewEncoder(transport).Encode(data))
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

func Test_Face_Put(t *testing.T) {
	transport := &bufferTransport{}
	subject := NewFace(transport)

	packet := &packets.Data{}
	packet.SetName(name.New(name.Component{"a"}))
	packet.SetContent([]byte("hello world"))

	// Publish a data packet
	err := subject.Put(packet)
	assert.NoError(t, err)

	// Check some data was written
	assert.True(t, len(transport.Buffer.Bytes()) > 0)
}

func Test_Face_RegisterPrefix(t *testing.T) {
	transport := &bufferTransport{}
	subject := NewFace(transport)
	keyChain := NewKeyChain()
	subject.SetCommandSigningInfo(keyChain, name.New(name.Component{"certificate"}))

	interests, err := subject.RegisterPrefix(name.New(name.Component{"a"}))
	assert.NoError(t, err)
	assert.NotNil(t, interests)

	// TODO: Need to send a packet through and check it's passed out on interests chan

	// Process the data, EOF is silenced
	// Do we need this?
	//assert.NoError(t, subject.ProcessEvents())

	// Check data was written
	assert.Equal(t, transport.Buffer.Bytes(), []byte("something here"))
}

func Test_Face_RegisterPrefix_NoCommandKeychainSet(t *testing.T) {
	transport := &bufferTransport{}
	subject := NewFace(transport)

	subject.SetCommandSigningInfo(nil, name.New(name.Component{"certificate"}))

	interests, err := subject.RegisterPrefix(name.New(name.Component{"a"}))
	assert.Equal(t, err, ErrCommandKeyChainNotSet)
	assert.Nil(t, interests)

	// Check no data was written
	assert.Equal(t, len(transport.Buffer.Bytes()), 0)
}

func Test_Face_RegisterPrefix_NoCommandCertificateSet(t *testing.T) {
	transport := &bufferTransport{}
	subject := NewFace(transport)
	keyChain := NewKeyChain()
	subject.SetCommandSigningInfo(keyChain, name.New())

	interests, err := subject.RegisterPrefix(name.New(name.Component{"a"}))
	assert.Equal(t, err, ErrCommandCertificateNameNotSet)
	assert.Nil(t, interests)

	// Check no data was written
	assert.Equal(t, len(transport.Buffer.Bytes()), 0)
}

func Test_Face_RegisterPrefix_ErrorRegistering(t *testing.T) {
	t.Error("pending")
}
