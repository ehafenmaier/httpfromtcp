package server

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode `json:"statusCode"`
	Message    string              `json:"message"`
}

type Handler func(w *response.Writer, req *request.Request)
