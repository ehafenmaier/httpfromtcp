package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
)

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
