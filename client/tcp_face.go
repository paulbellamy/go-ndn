package client

import "net"

func TCPFace(addr string) (*face, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return Face(conn), nil
}
