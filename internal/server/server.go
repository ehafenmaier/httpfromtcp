package server

import (
	"fmt"
	"httpfromtcp/internal/response"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("Error opening TCP listener on port %d: %v\n", port, err)
		return nil, err
	}

	server := &Server{listener: listener}
	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}

			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	err := response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		fmt.Printf("Error writing status line: %v\n", err)
	}

	err = response.WriteHeaders(conn, response.GetDefaultHeaders(0))
	if err != nil {
		fmt.Printf("Error writing headers: %v\n", err)
	}
}
