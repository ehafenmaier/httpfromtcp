package request

import (
	"fmt"
	"io"
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
	// Read the request line from the reader
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	rl, err := parseRequestLine(b)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request line: %w", err)
	}

	return &Request{
		RequestLine: rl,
	}, nil
}

func parseRequestLine(b []byte) (RequestLine, error) {
	// Split the request by CRLF and grab the first line (request line)
	req := string(b)
	rl := strings.Split(req, "\r\n")[0]

	// Split the request line into parts
	parts := strings.Split(rl, " ")

	// Verify we have exactly 3 parts: Method, RequestTarget, and HTTP Version
	if len(parts) != 3 {
		return RequestLine{}, fmt.Errorf("invalid request line: %s", rl)
	}

	// Verify the method part only contains capital letters
	for _, r := range parts[0] {
		if r < 'A' || r > 'Z' {
			return RequestLine{}, fmt.Errorf("invalid method in request line: %s", parts[0])
		}
	}

	// Verify that we have HTTP version 1.1
	if strings.TrimPrefix(parts[2], "HTTP/") != "1.1" {
		return RequestLine{}, fmt.Errorf("invalid HTTP version in request line: %s", parts[2])
	}

	return RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   strings.TrimPrefix(parts[2], "HTTP/"),
	}, nil
}
