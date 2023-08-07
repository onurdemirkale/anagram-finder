package api

import (
	"log"
	"net/http"
)

type HTTPError struct {
	Code    int
	Message string
}

const (
	ErrProcessing             = "an error occurred while processing your request"
	ErrInvalidFormat          = "invalid request format"
	ErrInvalidFile            = "failed to read file"
	ErrUnsupportedContentType = "unsupported content type"
	ErrInvalidInput           = "invalid input provided"
	ErrInvalidInputType       = "invalid input type. supported types: http, http_file, http_url"
	ErrInvalidAlgorithmType   = "invalid algorithm type. supported algorithms: basic"
	ErrInvalidFileInput       = "input data should be empty for file input type"
)

var ErrorMapping = map[string]HTTPError{
	ErrInvalidInput:           {http.StatusBadRequest, ErrInvalidInput},
	ErrInvalidInputType:       {http.StatusBadRequest, ErrInvalidInputType},
	ErrInvalidAlgorithmType:   {http.StatusBadRequest, ErrInvalidAlgorithmType},
	ErrInvalidFileInput:       {http.StatusBadRequest, ErrInvalidFileInput},
	ErrInvalidFormat:          {http.StatusBadRequest, ErrInvalidFormat},
	ErrInvalidFile:            {http.StatusBadRequest, ErrInvalidFile},
	ErrUnsupportedContentType: {http.StatusBadRequest, ErrUnsupportedContentType},
}

func handleError(err error) (int, string) {
	log.Printf("handler error: %v", err)
	if httpErr, ok := ErrorMapping[err.Error()]; ok {
		return httpErr.Code, httpErr.Message
	}
	return http.StatusInternalServerError, ErrProcessing
}
