package api

import (
	"errors"
	"strings"
)

const (
	inputTypeFile  = "http_file"
	inputTypeData  = "http"
	algorithmBasic = "basic"
)

type AnagramRequest struct {
	InputType string `json:"inputType"`
	InputData string `json:"inputData"`
	Algorithm string `json:"algorithm"`
}

func (req *AnagramRequest) validate() error {
	if err := req.validateInputType(); err != nil {
		return err
	}

	if err := req.validateAlgorithm(); err != nil {
		return err
	}

	if err := req.validateInputData(); err != nil {
		return err
	}

	return nil
}

func (req *AnagramRequest) validateInputType() error {
	supportedTypes := map[string]bool{
		inputTypeData: true,
		inputTypeFile: true,
	}

	if !supportedTypes[req.InputType] {
		return errors.New(ErrInvalidInputType)
	}

	return nil
}

func (req *AnagramRequest) validateAlgorithm() error {
	supportedAlgorithms := map[string]bool{
		algorithmBasic: true,
	}

	if !supportedAlgorithms[req.Algorithm] {
		return errors.New(ErrInvalidAlgorithmType)
	}

	return nil
}

func (req *AnagramRequest) validateInputData() error {
	// todo: improve file and http body input validations
	switch req.InputType {
	case inputTypeData:
		if req.InputData == "" || len(strings.Split(req.InputData, ",")) < 2 {
			return errors.New(ErrInvalidInput)
		}
	case inputTypeFile:
		if req.InputData != "" {
			return errors.New(ErrInvalidFileInput)
		}
	}

	return nil
}
