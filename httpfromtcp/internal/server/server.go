package server

import (
	"bytes"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
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

type Handler func(w io.Writer, req *request.Request) *HandlerError

type Server struct {
	listener net.Listener
	IsOpen   atomic.Bool
}

func (he *HandlerError) WriteError(w io.Writer) error {
	err := response.WriteStatusLine(w, he.StatusCode)
	if err != nil {
		return err
	}
	h := response.GetDefaultHeaders(len(he.Message))
	response.WriteHeaders(w, h)
	_, err = w.Write([]byte(he.Message))
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
			log.Fatal("error", err)
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
	body := new(bytes.Buffer)
	herror := handler(body, req)
	if herror != nil {
		err := herror.WriteError(conn)
		if err != nil {
			log.Println("Error writing handler error to the connection")
		}
	} else {
		res := new(bytes.Buffer)
		err = response.WriteStatusLine(res, 200)
		if err != nil {
			intError := HandlerError{StatusCode: 500, Message: "Error Writing Status Line"}
			err := intError.WriteError(conn)
			if err != nil {
				log.Println("Error Writing Internal Error to the Connection")
			}
			return
		}
		h := response.GetDefaultHeaders(body.Len())
		err = response.WriteHeaders(res, h)
		if err != nil {
			intError := HandlerError{StatusCode: 500, Message: "Error Writing Headers"}
			err := intError.WriteError(conn)
			if err != nil {
				log.Println("Error Writing Internal Error to the Connection")
			}
			return
		}
		_, err = body.WriteTo(res)
		if err != nil {
			intError := HandlerError{StatusCode: 500, Message: "Error Writing Response Body"}
			err := intError.WriteError(conn)
			if err != nil {
				log.Println("Error Writing Internal Error to the Connection")
			}
			return
		}
		_, err := res.WriteTo(conn)
		if err != nil {
			intError := HandlerError{StatusCode: 500, Message: "Error Writing Response to Connection"}
			err := intError.WriteError(conn)
			if err != nil {
				log.Println("Error Writing Internal Error to the Connection")
			}
			return
		}
	}
}
