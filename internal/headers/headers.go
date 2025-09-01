package headers

import (
	"bytes"
	"fmt"
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
		line := string(data[bytesConsumed : bytesConsumed+idx])
		parts := strings.SplitN(line, ":", 2)

		// If there is any whitespace at the end of the key the header line is invalid
		if strings.TrimRight(parts[0], " ") != parts[0] {
			return 0, false, fmt.Errorf("invalid header line: %s", line)
		}

		// Add the key value pair to the headers map
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		h[key] = value

		// Increment the number of bytes consumed
		bytesConsumed += idx + len(rn)
	}
}
