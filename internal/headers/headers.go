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

	// Find the index of the first CRLF
	idx := bytes.Index(data, rn)

	// If not found we need more data
	if idx == -1 {
		return 0, false, nil
	}

	// If we find a CRLF at the start, we're done with headers
	// Return the consumed CRLF bytes
	if idx == 0 {
		return len(rn), true, nil
	}

	// Parse the header line
	line := string(data[:idx])
	key, value, err := parseHeaderLine(line)
	if err != nil {
		return 0, false, err
	}

	// Check if the header key exists and if it does add the value to a
	// comma-delimited list otherwise add the key and value
	v, exists := h[key]
	if exists {
		h[key] = v + ", " + value
	} else {
		h[key] = value
	}

	// Return the number of bytes consumed plus the CRLF
	return idx + len(rn), false, nil
}

func (h Headers) Get(key string) (string, bool) {

	v, ok := h[strings.ToLower(key)]
	return v, ok
}

func parseHeaderLine(line string) (string, string, error) {
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

	// Header line is valid so return lowercase key and value
	key := strings.ToLower(strings.TrimSpace(parts[0]))
	value := strings.TrimSpace(parts[1])

	return key, value, nil
}
