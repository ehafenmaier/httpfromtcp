package request

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strings"
)

const bufferSize = 8

type RequestState int

const (
	Initialized RequestState = iota
	ParsingHeaders
	Done
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	State       RequestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	// Create a buffer to read data into
	buf := make([]byte, bufferSize)
	readToIndex := 0

	// Create a new Request in the Initialized state
	req := &Request{
		State:   Initialized,
		Headers: map[string]string{},
	}

	for req.State != Done {
		// If the buffer is full double its size
		if readToIndex == len(buf) {
			doubleBuf := make([]byte, len(buf)*2)
			copy(doubleBuf, buf)
			buf = doubleBuf
		}

		// Read data from the reader into the buffer
		br, err := reader.Read(buf[readToIndex:])
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("error reading from reader: %w", err)
		}

		// If we read 0 bytes and got EOF, break
		if br == 0 && err == io.EOF {
			req.State = Done
			continue
		}

		// Parse the data in the buffer
		readToIndex += br
		bp, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, fmt.Errorf("error parsing request: %w", err)
		}

		// If no bytes were parsed, we need more data
		if bp == 0 {
			continue
		}

		// Remove the parsed bytes from the buffer
		newBuf := make([]byte, len(buf))
		copy(newBuf, buf[bp:])
		buf = newBuf
		readToIndex -= bp
	}

	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesConsumed := 0
	for r.State != Done {
		bc, err := r.parseSingle(data[totalBytesConsumed:])
		// If there's an error, return it
		if err != nil {
			return 0, fmt.Errorf("error parsing request: %w", err)
		}

		totalBytesConsumed += bc

		// If no bytes were consumed, we need more data
		if totalBytesConsumed == 0 {
			return 0, nil
		} else {
			break
		}
	}

	return totalBytesConsumed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.State {
	case Initialized:
		// If the request is initialized, parse the request line
		rl, bc, err := parseRequestLine(data)
		// If there was an error parsing, return it
		if err != nil {
			return 0, err
		}
		// If no bytes were consumed, we need more data
		if bc == 0 {
			return 0, nil
		}
		// Successfully parsed the request line
		r.RequestLine = rl
		r.State = ParsingHeaders
		return bc, nil
	case ParsingHeaders:
		// Parse the headers
		bc, done, err := r.Headers.Parse(data)
		// If there's an error parsing, return it
		if err != nil {
			return 0, err
		}
		// If headers are completely parsed, set the request state to done
		if done {
			r.State = Done
		}
		// If no bytes were consumed, we need more data
		if bc == 0 {
			return 0, nil
		}
		// Return the total bytes consumed
		return bc, nil
	case Done:
		// If the request is done, something went wrong
		return 0, fmt.Errorf("trying to read data in a done state")
	default:
		return 0, fmt.Errorf("unknown request state")
	}
}

func parseRequestLine(b []byte) (RequestLine, int, error) {
	// Find the index of the first CRLF to isolate the request line
	rn := []byte("\r\n")
	idx := bytes.Index(b, rn)
	if idx == -1 {
		return RequestLine{}, 0, nil
	}

	bytesConsumed := len(b)
	rl := string(b[:idx])

	// Split the request line into parts
	parts := strings.Split(rl, " ")

	// Verify we have exactly 3 parts: Method, RequestTarget, and HTTP Version
	if len(parts) != 3 {
		return RequestLine{}, 0, fmt.Errorf("invalid request line: %s", rl)
	}

	// Verify the method part only contains capital letters
	for _, r := range parts[0] {
		if r < 'A' || r > 'Z' {
			return RequestLine{}, 0, fmt.Errorf("invalid method in request line: %s", parts[0])
		}
	}

	// Verify that we have HTTP version 1.1
	if strings.TrimPrefix(parts[2], "HTTP/") != "1.1" {
		return RequestLine{}, 0, fmt.Errorf("invalid HTTP version in request line: %s", parts[2])
	}

	return RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   strings.TrimPrefix(parts[2], "HTTP/"),
	}, bytesConsumed, nil
}
