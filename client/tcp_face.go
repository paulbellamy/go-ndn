package client

func TCPFace(addr string) *face {
	return Face(TCPTransport(addr))
}
