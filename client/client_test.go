package client

import (
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
