package server

import "net"

type Server struct {
}

func (s *Server) Addr() net.Addr {
	return nil
}

func (s *Server) Close() error {
	return nil
}
