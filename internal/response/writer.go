package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
)

type Writer struct {
	writer      io.Writer
	writerState WriterState
}

type WriterState int

const (
	StatusLine WriterState = iota
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
		return fmt.Errorf("writer state is not StatusLine")
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
		return fmt.Errorf("writer state is not Headers")
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
		return 0, fmt.Errorf("writer state is not Body")
	}

	return w.writer.Write(p)
}
