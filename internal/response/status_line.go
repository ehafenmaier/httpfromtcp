package response

import "io"

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
