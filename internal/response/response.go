package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
)

type StatusCode uint

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest                     = 400
	StatusInternalServerError            = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
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

	_, err := w.Write(response)
	if err != nil {
		return err
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.Headers{}

	headers["Content-Length"] = strconv.Itoa(contentLen)
	headers["Connection"] = "close"
	headers["Content-Type"] = "text/plain"

	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	var response []byte

	for key, value := range headers {
		response = append(response, []byte(fmt.Sprintf("%s: %s\r\n", key, value))...)
	}

	response = append(response, []byte("\r\n")...)

	_, err := w.Write(response)
	if err != nil {
		return err
	}

	return nil
}
