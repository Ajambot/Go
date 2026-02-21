package server

import (
	"httpfromtcp/pkg/request"
	"httpfromtcp/pkg/response"
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	listener net.Listener
	IsOpen   atomic.Bool
}

func (he *HandlerError) WriteError(w io.Writer) error {
	buf := new(response.Writer)
	err := buf.WriteStatusLine(he.StatusCode)
	if err != nil {
		return err
	}
	h := response.GetDefaultHeaders(len(he.Message))
	err = buf.WriteHeaders(h)
	if err != nil {
		return err
	}
	_, err = buf.WriteBody([]byte(he.Message))
	return err
}

func Serve(port int, handler Handler) (*Server, error) {
	lstnr, e := net.Listen("tcp", ":"+strconv.Itoa(port))
	if e != nil {
		return nil, e
	}
	server := Server{listener: lstnr}
	server.IsOpen.Store(true)
	go server.listen(handler)
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

func (s *Server) listen(handler Handler) {
	for {
		conn, err := s.listener.Accept()
		if !s.IsOpen.Load() {
			break
		}
		if err != nil {
			log.Print(err)
			continue
		}

		go s.handle(conn, handler)
	}
}

func (s *Server) handle(conn net.Conn, handler Handler) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		intError := HandlerError{StatusCode: 400, Message: err.Error()}
		err := intError.WriteError(conn)
		if err != nil {
			log.Println("Error Writing Internal Error to the Connection")
		}
		return
	}
	body := response.Writer{Buffer: conn}
	handler(&body, req)
}
