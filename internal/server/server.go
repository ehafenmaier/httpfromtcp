package server

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("Error opening TCP listener on port %d: %v\n", port, err)
		return nil, err
	}

	server := &Server{
		listener: listener,
		handler:  handler,
	}
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

	headers := response.GetDefaultHeaders(0)
	r, err := request.RequestFromReader(conn)
	if err != nil {
		response.WriteStatusLine(conn, response.StatusBadRequest)
		response.WriteHeaders(conn, headers)
		return
	}

	writer := bytes.NewBuffer([]byte{})
	handlerError := s.handler(writer, r)
	if handlerError != nil {
		response.WriteStatusLine(conn, handlerError.StatusCode)
		response.WriteHeaders(conn, headers)
		conn.Write([]byte(handlerError.Message))
		return
	}

	body := writer.Bytes()
	headers = response.GetDefaultHeaders(len(body))
	response.WriteStatusLine(conn, response.StatusOK)
	response.WriteHeaders(conn, headers)
	conn.Write(body)
}
