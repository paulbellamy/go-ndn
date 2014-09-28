package ndn

import "net"

func NewTCPFace(addr string) (*Face, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewFace(conn), nil
}
