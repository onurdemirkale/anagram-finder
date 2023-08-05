package inputsource

import (
	"strings"
)

type HttpBodyInputSource struct {
	words []string
}

func NewHttpBodyInputSource(inputData string) *HttpBodyInputSource {
	words := strings.Split(inputData, ",")
	return &HttpBodyInputSource{words: words}
}

func (h *HttpBodyInputSource) GetWords() ([]string, error) {
	return h.words, nil
}
