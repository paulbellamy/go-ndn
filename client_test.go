package ndn

import (
	"bytes"
	"io"
	"testing"

	"github.com/paulbellamy/go-ndn/server"
	"github.com/stretchr/testify/mock"
)

func testServer(t *testing.T) *server.Server {
	return &server.Server{}
}

func testTransport(t *testing.T) Transport {
	return (Transport)(nil)
}

type mockTransport struct {
	mock.Mock
}

func (m *mockTransport) Write(b []byte) (int, error) {
	args := m.Mock.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *mockTransport) Read(b []byte) (int, error) {
	args := m.Mock.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *mockTransport) Close() error {
	return m.Mock.Called().Error(0)
}

type bufferTransport struct {
	bytes.Buffer
	closed    bool
	remaining io.Reader
}

func (b *bufferTransport) Read(p []byte) (int, error) {
	if b.closed {
		return b.remaining.Read(p)
	}
	return b.Buffer.Read(p)
}

func (b *bufferTransport) Close() error {
	b.closed = true
	b.remaining = bytes.NewReader(b.Buffer.Bytes())
	return nil
}
