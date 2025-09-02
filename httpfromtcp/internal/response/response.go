package response

import (
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
)

type StatusCode int

const (
	WritingStatusLine int = iota
	WritingHeaders
	WritingBody
	Done
)

type Writer struct {
	Buffer      io.Writer
	writerState int
}

const (
	OK                    StatusCode = 200
	BAD_REQUEST           StatusCode = 400
	INTERNAL_SERVER_ERROR StatusCode = 500
)

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != WritingStatusLine {
		return errors.New("Error: Writer not in WritingStatusLine state")
	}
	switch statusCode {
	case OK:
		_, err := w.Buffer.Write([]byte("HTTP/1.1 200 OK\r\n"))
		w.writerState = WritingHeaders
		return err
	case BAD_REQUEST:
		_, err := w.Buffer.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		w.writerState = WritingHeaders
		return err
	case INTERNAL_SERVER_ERROR:
		_, err := w.Buffer.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		w.writerState = WritingHeaders
		return err
	default:
		_, err := w.Buffer.Write([]byte("HTTP/1.1 " + fmt.Sprint(statusCode) + " \r\n"))
		w.writerState = WritingHeaders
		return err
	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprint(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != WritingHeaders {
		return errors.New("Error: Writer not in WritingHeaders state")
	}

	ch := headers.Range()
	for h := range ch {
		_, err := w.Buffer.Write([]byte(h[0] + ": " + h[1] + "\r\n"))
		if err != nil {
			return err
		}
	}
	_, err := w.Buffer.Write([]byte("\r\n"))
	w.writerState = WritingBody
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != WritingBody {
		return 0, errors.New("Error: Writer not in WritingBody state")
	}

	n, err := w.Buffer.Write(p)
	w.writerState = Done
	return n, err
}

func (w *Writer) Write(p []byte) (n int, err error) {
	n, err = w.Buffer.Write(p)
	return n, err
}
