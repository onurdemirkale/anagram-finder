package inputsource

import (
	"fmt"
	"mime/multipart"
)

type InputSourceFactoryInterface interface {
	CreateInputSource(inputType string, inputData interface{}) (InputSource, error)
}

type InputSourceFactory struct{}

func NewInputSourceFactory() InputSourceFactoryInterface {
	return &InputSourceFactory{}
}

func (f *InputSourceFactory) CreateInputSource(inputType string, inputData interface{}) (InputSource, error) {
	switch inputType {
	case "http_body":
		return NewHttpBodyInputSource(inputData.(string)), nil
	case "http_file":
		// todo: add error handling for type assertion to avoid potential panics
		return NewHttpFileInputSource(inputData.(multipart.File)), nil
	case "http_url":
		return NewHttpUrlInputSource(inputData.(string)), nil
	default:
		return nil, fmt.Errorf("unknown input type: %s", inputType)
	}
}
