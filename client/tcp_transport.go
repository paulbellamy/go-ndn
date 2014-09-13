package client

import "net"

func TCPTransport(addr string) *tcpTransport {
	return &tcpTransport{}
}

type tcpTransport struct {
	net.Conn
}
