package server

import (
	"fmt"
	"net"
	"sync/atomic"
)

type Server struct {
	closed atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("Error opening TCP listener on port %d: %v\n", port, err)
		return nil, err
	}

	server := &Server{}
	go server.listen(listener)

	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return nil
}

func (s *Server) listen(listener net.Listener) {
	for {
		conn, err := listener.Accept()

		if s.closed.Load() {
			return
		}

		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	response := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!")
	_, err := conn.Write(response)
	if err != nil {
		fmt.Printf("Error writing message: %v\n", err)
	}
}
