package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	rn := []byte("\r\n")
	bytesConsumed := 0
	for {
		// Find the index of the first CRLF
		idx := bytes.Index(data[bytesConsumed:], rn)

		// If not found we need more data
		if idx == -1 {
			return 0, false, nil
		}

		// If we find a CRLF at the start, we're done with headers
		if idx == 0 {
			return bytesConsumed, true, nil
		}

		// Parse the header line
		key, value, err := parseHeaderLine(data[bytesConsumed : bytesConsumed+idx])
		if err != nil {
			return 0, false, err
		}

		// Add the key value pair to the headers map
		h[key] = value

		// Increment the number of bytes consumed
		bytesConsumed += idx + len(rn)
	}
}

func parseHeaderLine(data []byte) (string, string, error) {
	line := string(data)
	parts := strings.SplitN(line, ":", 2)

	// If we don't have exactly two parts the header line is invalid
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid header line: %s", line)
	}

	// If there is any whitespace at the end of the key the header line is invalid
	if strings.TrimRight(parts[0], " ") != parts[0] {
		return "", "", fmt.Errorf("invalid header line: %s", line)
	}

	// If the header key doesn't only have alphanumeric and certain special
	// characters it is invalid
	pattern := "^[a-zA-Z0-9!#$%&'*+-.^_`|~]+$"
	re := regexp.MustCompile(pattern)
	if !re.MatchString(strings.TrimSpace(parts[0])) {
		return "", "", fmt.Errorf("invalid header key: %s", parts[0])
	}

	// Header line is valid so return lowercased key and value
	key := strings.ToLower(strings.TrimSpace(parts[0]))
	value := strings.TrimSpace(parts[1])

	return key, value, nil
}
