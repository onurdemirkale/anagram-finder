package api

import (
	"errors"
	"fmt"
	"strings"
)

const (
	inputTypeData  = "http"
	algorithmBasic = "basic"
)

var (
	supportedInputTypes     = []string{"http"}
	supportedAlgorithms     = []string{"basic"}
	errInvalidInput         = errors.New("invalid input provided")
	errInvalidInputType     = fmt.Errorf("invalid input type. supported types: %s", strings.Join(supportedInputTypes, ", "))
	errInvalidAlgorithmType = fmt.Errorf("invalid algorithm type. supported algorithms: %s", strings.Join(supportedAlgorithms, ", "))
)

type AnagramRequest struct {
	InputType string `json:"inputType"`
	InputData string `json:"inputData"`
	Algorithm string `json:"algorithm"`
}

func (req *AnagramRequest) Validate() error {
	supportedAlgorithms := map[string]bool{
		algorithmBasic: true,
	}

	supportedTypes := map[string]bool{
		inputTypeData: true,
	}

	if !supportedTypes[req.InputType] {
		return errInvalidInputType
	}

	if !supportedAlgorithms[req.Algorithm] {
		return errInvalidAlgorithmType
	}

	switch req.InputType {
	case inputTypeData:
		if !isValidData(req.InputData) {
			return errInvalidInput
		}
	}

	return nil
}

// used for validating data in http body
func isValidData(data string) bool {
	return strings.Contains(data, ",")
}
