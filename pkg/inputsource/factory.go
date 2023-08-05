package inputsource

import "fmt"

type InputSourceFactoryInterface interface {
	CreateInputSource(inputType string, inputData string) (InputSource, error)
}

type InputSourceFactory struct{}

func NewInputSourceFactory() InputSourceFactoryInterface {
	return &InputSourceFactory{}
}

func (f *InputSourceFactory) CreateInputSource(inputType string, inputData string) (InputSource, error) {
	switch inputType {
	case "http":
		return NewHttpBodyInputSource(inputData), nil
	default:
		return nil, fmt.Errorf("unknown input type: %s", inputType)
	}
}
