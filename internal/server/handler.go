package server

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
)

type HandlerError struct {
	StatusCode response.StatusCode `json:"statusCode"`
	Message    string              `json:"message"`
}

type Handler func(w io.Writer, req *request.Request) *HandlerError
