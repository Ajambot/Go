package request

import (
	"errors"
	"io"
	"strings"
)

const bufSize int = 8

type Request struct {
	RequestLine RequestLine
	Status      int // 0 is initialized, 1 is done
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := &Request{Status: 0}

	buf := make([]byte, bufSize)
	readToIndex := 0

	for req.Status != 1 {
		if len(buf) == cap(buf) {
			newBuf := make([]byte, cap(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		n, err := reader.Read(buf[readToIndex:])
		if err == io.EOF {
			req.Status = 1
			break
		}

		readToIndex += n

		n, err = req.parse(buf)
		if err != nil {
			return nil, err
		}
		if n > 0 {
			copy(buf, buf[readToIndex:])
			readToIndex -= n
		}
	}
	return req, nil
}

func parseRequestLine(request string) (*RequestLine, int, error) {
	lines := strings.Split(request, "\r\n")
	if len(lines) == 1 {
		return nil, 0, nil
	}
	reqLine := string(lines[0])
	parts := strings.Split(reqLine, " ")

	if len(parts) < 3 {
		return nil, 0, errors.New("Invalid number of parts in request line")
	}

	method := parts[0]
	target := parts[1]
	version := parts[2]

	for _, c := range method {
		if !((65 <= c) && (c <= 90)) {
			return nil, 0, errors.New("Request method is not valid")
		}
	}

	if version[:5] != "HTTP/" {
		return nil, 0, errors.New("Version section not valid.")
	}

	if version[5:] != "1.1" {
		return nil, 0, errors.New("Only HTTP version 1.1 supported")
	}

	version = version[5:]

	return &RequestLine{HttpVersion: version, RequestTarget: target, Method: method}, len(reqLine), nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.Status == 1 {
		return 0, errors.New("Error: trying to read data in a done state")
	} else if r.Status != 0 {
		return 0, errors.New("Error: unknown state")
	}

	reqLine, n, err := parseRequestLine(string(data))

	if n > 0 {
		r.RequestLine = *reqLine
		r.Status = 1
	}

	return n, err
}
