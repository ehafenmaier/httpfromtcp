package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strings"
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
	w.writer.Write([]byte(fmt.Sprintf("%04X\r\n", len(p))))
	// Append carriage return to data and write
	p = append(p, []byte("\r\n")...)
	w.writer.Write(p)

	return len(p), nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	w.writer.Write([]byte("0\r\n\r\n"))

	return 0, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	w.writer.Write([]byte("0\r\n"))
	trailers, ok := h.Get("Trailer")
	if !ok {
		return fmt.Errorf("trailer header not found")
	}

	trailerKeys := strings.Split(trailers, ", ")
	for _, key := range trailerKeys {
		value, ok := h.Get(key)
		if ok {
			w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
		}
	}
	w.writer.Write([]byte("\r\n"))

	return nil
}
