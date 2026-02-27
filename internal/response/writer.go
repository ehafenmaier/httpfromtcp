package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
)

type Writer struct {
	writer      io.Writer
	writerState writerState
}

type writerState int

const (
	StatusLine writerState = iota
	Headers
	Body
)

const CRLF = "\r\n"

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer:      w,
		writerState: StatusLine,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != StatusLine {
		return fmt.Errorf("cannot write status line in state %d", w.writerState)
	}

	var response []byte

	switch statusCode {
	case StatusOK:
		response = []byte("HTTP/1.1 200 OK\r\n")
	case StatusBadRequest:
		response = []byte("HTTP/1.1 400 Bad Request\r\n")
	case StatusInternalServerError:
		response = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	default:
		response = []byte("HTTP/1.1 400 Bad Request\r\n")
	}

	_, err := w.writer.Write(response)
	if err != nil {
		return err
	}

	w.writerState = Headers
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != Headers {
		return fmt.Errorf("cannot write headers in state %d", w.writerState)
	}

	var response []byte

	for key, value := range headers {
		response = append(response, []byte(fmt.Sprintf("%s: %s\r\n", key, value))...)
	}

	response = append(response, []byte("\r\n")...)

	_, err := w.writer.Write(response)
	if err != nil {
		return err
	}

	w.writerState = Body
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != Body {
		return 0, fmt.Errorf("cannot write body in state %d", w.writerState)
	}

	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.writerState != Body {
		return 0, fmt.Errorf("cannot write chunked body in state %d", w.writerState)
	}

	// Write length of data in hex
	w.writer.Write([]byte(fmt.Sprintf("%04X%s", len(p), CRLF)))
	// Append carriage return to data and write
	p = append(p, []byte(CRLF)...)
	w.writer.Write(p)

	return len(p), nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	w.writer.Write([]byte("0" + CRLF))
	w.writer.Write([]byte(CRLF))

	return 0, nil
}
