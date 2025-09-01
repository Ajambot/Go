package request

import (
	"errors"
	"httpfromtcp/internal/headers"
	"io"
	"strings"
)

const bufSize int = 8

type Request struct {
	RequestLine RequestLine
	Status      int
	Headers     headers.Headers
}

const (
	Initialized    int = iota // 0
	ParsingHeaders            // 1
	Done                      // 2
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := &Request{Headers: headers.Headers{}, Status: Initialized}

	buf := make([]byte, bufSize)
	readToIndex := 0

	for req.Status != Done {
		if readToIndex >= cap(buf) {
			newBuf := make([]byte, cap(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		n, err := reader.Read(buf[readToIndex:])
		if err == io.EOF {
			return nil, errors.New("Error: Received End of Input while Parser Status Is Not Done")
		}

		readToIndex += n

		n, err = req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		if n > 0 {
			copy(buf, buf[n:])
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

	return &RequestLine{HttpVersion: version, RequestTarget: target, Method: method}, len(reqLine) + len("\r\n"), nil
}

func (r *Request) parseSingleHeader(data []byte) (int, error) {
	n, done, err := r.Headers.Parse(data)

	if err != nil {
		return 0, err
	}

	if done {
		r.Status = Done
	}

	return n, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.Status {
	case Initialized:
		reqLine, n, err := parseRequestLine(string(data))

		if n > 0 {
			r.RequestLine = *reqLine
			r.Status = ParsingHeaders
		}

		return n, err

	case ParsingHeaders:
		totalBytesParsed := 0
		for {
			n, err := r.parseSingleHeader(data[totalBytesParsed:])
			totalBytesParsed += n
			if err != nil || n == 0 || totalBytesParsed >= len(data) {
				return totalBytesParsed, err
			}
		}

	case Done:
		return 0, errors.New("Error: trying to read data in a done state")
	default:
		return 0, errors.New("Error: Unexpected Parser State")
	}
}
