package headers

import (
	"errors"
	"strings"
)

type Headers struct {
	values map[string]string
}

func NewHeaders() Headers {
	return Headers{values: make(map[string]string)}
}

func (h Headers) Get(key string) (string, bool) {
	v, ok := h.values[strings.ToLower(key)]

	return v, ok
}

func (h Headers) Set(key, value string) {
	if v, ok := h.values[strings.ToLower(key)]; ok {
		h.values[strings.ToLower(key)] = v + ", " + value
	} else {
		h.values[strings.ToLower(key)] = value
	}
}

func (h Headers) Length() int {
	return len(h.values)
}

func (h Headers) Range() <-chan [2]string {
	keyValChan := make(chan [2]string)
	go func() {
		for k, v := range h.values {
			var keyVal = [2]string{k, v}
			keyValChan <- keyVal
		}
		close(keyValChan)
	}()
	return keyValChan
}

const crlf string = "\r\n"
const validChars string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$%&'*+-.^_`|~"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	str := string(data)
	if !strings.Contains(str, crlf) {
		return 0, false, nil
	}
	crlfIndex := strings.Index(str, crlf)
	if crlfIndex == 0 {
		return len(data[:crlfIndex+len(crlf)]), true, nil
	}

	str = str[:crlfIndex]
	str = strings.TrimSpace(str)

	sepIndex := strings.Index(str, ":")
	if sepIndex == -1 {
		return 0, false, errors.New("Headers malformed. Missing colon")
	}

	fieldName := str[:sepIndex]
	fieldValue := str[sepIndex+1:]

	for _, c := range fieldName {
		if !strings.Contains(validChars, string(c)) {
			return 0, false, errors.New("Headers malformed. Non-valid character in header")
		}
	}

	fieldName = strings.TrimSpace(fieldName)
	fieldValue = strings.TrimSpace(fieldValue)

	h.Set(fieldName, fieldValue)

	return len(data[:crlfIndex+len(crlf)]), false, nil
}
