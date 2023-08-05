package inputsource

import "fmt"

type InputSourceFactory struct{}

func NewInputSourceFactory() *InputSourceFactory {
	return &InputSourceFactory{}
}

func (f *InputSourceFactory) CreateInputSource(inputType string, inputData string) (InputSource, error) {

	fmt.Println(inputType)
	switch inputType {
	case "http":
		return NewHttpBodyInputSource(inputData), nil
	default:
		return nil, fmt.Errorf("unknown input type: %s", inputType)
	}
}
