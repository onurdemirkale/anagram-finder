package api

import (
	"errors"
	"fmt"
	"strings"
)

const (
	inputTypeFile  = "http_file"
	inputTypeData  = "http"
	algorithmBasic = "basic"
)

var (
	supportedInputTypes     = []string{"http", "http_file"}
	supportedAlgorithms     = []string{"basic"}
	errInvalidInput         = errors.New("invalid input provided")
	errInvalidInputType     = fmt.Errorf("invalid input type. supported types: %s", strings.Join(supportedInputTypes, ", "))
	errInvalidAlgorithmType = fmt.Errorf("invalid algorithm type. supported algorithms: %s", strings.Join(supportedAlgorithms, ", "))
	errInvalidFileInput     = errors.New("input data should be empty for file input type")
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
		inputTypeFile: true,
	}

	if !supportedTypes[req.InputType] {
		return errInvalidInputType
	}

	if !supportedAlgorithms[req.Algorithm] {
		return errInvalidAlgorithmType
	}

	// todo: improve file and http body input validations
	switch req.InputType {
	case inputTypeData:
		if !strings.Contains(req.InputData, ",") {
			return errInvalidInput
		}
	case inputTypeFile:
		if req.InputData != "" {
			return errInvalidFileInput
		}
	}

	return nil
}
