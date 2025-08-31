package headers

import (
	"errors"
	"strings"
	"unicode"
)

type Headers map[string]string

var crlf []byte = []byte("\r\n")

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	str := string(data)
	if !strings.Contains(str, "\r\n") {
		return 0, false, nil
	}
	crlfIndex := strings.Index(str, string(crlf))
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
		if unicode.IsSpace(c) {
			return 0, false, errors.New("Headers malformed. Whitespace in field name")
		}
	}

	fieldName = strings.TrimSpace(fieldName)
	fieldValue = strings.TrimSpace(fieldValue)

	h[fieldName] = fieldValue
	return len(data[:crlfIndex+len(crlf)]), false, nil
}
