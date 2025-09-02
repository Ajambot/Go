package server

import (
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	IsOpen   atomic.Bool
}

func Serve(port int) (*Server, error) {
	lstnr, e := net.Listen("tcp", ":"+strconv.Itoa(port))
	server := Server{listener: lstnr}
	server.IsOpen.Store(true)
	if e != nil {
		return nil, e
	}
	go server.listen()
	return &server, nil
}

func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return err
	}
	s.IsOpen.Store(false)
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if !s.IsOpen.Load() {
			break
		}
		if err != nil {
			log.Fatal("error", err)
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	res := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!\n"

	conn.Write([]byte(res))
	conn.Close()
}
