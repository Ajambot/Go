package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

const crlf string = "\r\n"
const validChars string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$%&'*+-.^_`|~"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	str := string(data)
	if !strings.Contains(str, "\r\n") {
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

	h[strings.ToLower(fieldName)] = fieldValue
	return len(data[:crlfIndex+len(crlf)]), false, nil
}
