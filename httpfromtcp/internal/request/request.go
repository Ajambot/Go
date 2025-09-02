package request

import (
	"errors"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
	"strings"
)

const bufSize int = 8

type Request struct {
	RequestLine RequestLine
	Status      int
	Headers     headers.Headers
	Body        []byte
}

const (
	Initialized    int = iota // 0
	ParsingHeaders            // 1
	ParsingBody               // 2
	Done                      // 3
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := &Request{Headers: headers.NewHeaders(), Status: Initialized}

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
			return nil, errors.New("Reached EOF Before Finishing Reading Request")
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
		r.Status = ParsingBody
	}

	return n, nil
}

func (r *Request) parseBody(data []byte) (int, error) {
	v, ok := r.Headers.Get("Content-Length")

	if !ok {
		r.Status = Done
		return 0, nil
	}

	hLen, err := strconv.Atoi(v)

	if err != nil {
		return 0, errors.New("Content-Length is not a number")
	}

	if len(data) > hLen {
		return 0, errors.New("Body bigger than Content-Length")
	}

	if len(data) < hLen {
		return 0, nil
	}

	r.Body = make([]byte, len(data))
	copy(r.Body, data)
	r.Status = Done
	return len(data), nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for {
		switch r.Status {
		case Initialized:
			reqLine, n, err := parseRequestLine(string(data[totalBytesParsed:]))

			totalBytesParsed += n

			if err != nil || n == 0 {
				return totalBytesParsed, err
			}

			if n > 0 {
				r.RequestLine = *reqLine
				r.Status = ParsingHeaders
			}

		case ParsingHeaders:
			n, err := r.parseSingleHeader(data[totalBytesParsed:])
			totalBytesParsed += n
			if err != nil || n == 0 {
				return totalBytesParsed, err
			}

		case ParsingBody:
			n, err := r.parseBody(data[totalBytesParsed:])

			totalBytesParsed += n

			if err != nil || n == 0 || totalBytesParsed >= len(data) {
				return totalBytesParsed, err
			}

		case Done:
			return 0, errors.New("Error: trying to read data in a done state")
		default:
			return 0, errors.New("Error: Unexpected Parser State")
		}
	}

}
