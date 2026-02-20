package request

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
	"strings"
)

const bufferSize = 8

type RequestState int

const (
	Initialized RequestState = iota
	ParsingHeaders
	ParsingBody
	Done
)

type Request struct {
	RequestLine   RequestLine
	Headers       headers.Headers
	Body          []byte
	ContentLength int
	State         RequestState
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
		Body:    make([]byte, 0),
	}

	for req.State != Done {
		// If the buffer is full double its size
		if readToIndex >= len(buf) {
			doubleBuf := make([]byte, len(buf)*2)
			copy(doubleBuf, buf)
			buf = doubleBuf
		}

		// Read data from the reader into the buffer
		br, err := reader.Read(buf[readToIndex:])
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("error reading from reader: %w", err)
		}

		// If we read 0 bytes and got EOF we're done
		if br == 0 && err == io.EOF {
			// Check for valid content length
			if req.ContentLength > 0 && len(req.Body) < req.ContentLength {
				return nil, fmt.Errorf("request body smaller than specified content length")
			}
			req.State = Done
			continue
		}

		// Parse the data in the buffer
		readToIndex += br
		bp, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, fmt.Errorf("error parsing request: %w", err)
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
		if bc == 0 {
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
		// If headers are completely parsed set the request state to parsing body
		if done {
			r.State = ParsingBody
		}
		// Return the total bytes consumed
		return bc, nil
	case ParsingBody:
		// Check for content length header and if it doesn't exist we're done
		contentLen, ok := r.Headers.Get("Content-Length")
		if !ok {
			r.State = Done
			return 0, nil
		}
		// Convert content length header value to an integer
		cl, err := strconv.Atoi(contentLen)
		if err != nil {
			return 0, err
		}
		// Set the content length to the header value and the
		// body to the remaining data minus the leading "\r\n"
		r.ContentLength = cl
		r.Body = append(r.Body, data...)
		// If body is larger than specified content length return error
		if len(r.Body) > r.ContentLength {
			return 0, fmt.Errorf("request body larger than specified content length")
		}
		// If body equals specified content length we're done
		if len(r.Body) == r.ContentLength {
			r.State = Done
			return 0, nil
		}
		// We need more data
		return 0, nil
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

	bytesConsumed := idx + len(rn)
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
