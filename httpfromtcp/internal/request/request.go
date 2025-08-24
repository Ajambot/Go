package request

import (
	"errors"
	"io"
	"log"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r, err := io.ReadAll(reader)

	if err != nil {
		log.Fatal("error ", err)
	}

	reqLine, err := parseRequestLine(string(r))

	if err != nil {
		return nil, err
	}

	return &Request{RequestLine: *reqLine}, nil
}

func parseRequestLine(request string) (*RequestLine, error) {
	reqLine := strings.Split(request, "\r\n")[0]
	parts := strings.Split(reqLine, " ")

	if len(parts) < 3 {
		return nil, errors.New("Invalid number of parts in request line")
	}

	method := parts[0]
	target := parts[1]
	version := parts[2]

	for _, c := range method {
		if !((65 <= c) && (c <= 90)) {
			return nil, errors.New("Request method is not valid")
		}
	}

	if version[:5] != "HTTP/" {
		return nil, errors.New("Version section not valid.")
	}

	if version[5:] != "1.1" {
		return nil, errors.New("Only HTTP version 1.1 supported")
	}

	version = version[5:]

	return &RequestLine{HttpVersion: version, RequestTarget: target, Method: method}, nil
}
