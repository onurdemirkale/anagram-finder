package inputsource

import (
	"bufio"
	"net/http"
	"strings"
)

type HttpUrlInputSource struct {
	url string
}

const maxBufSize = 2 * 1024 * 1024 // 2MB

func NewHttpUrlInputSource(url string) *HttpUrlInputSource {
	return &HttpUrlInputSource{url: url}
}

// todo: add error checks for the response status codes
func (hu *HttpUrlInputSource) GetWords() ([]string, error) {
	resp, err := http.Get(hu.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var words []string

	scanner := bufio.NewScanner(resp.Body)
	buf := make([]byte, 0, bufio.MaxScanTokenSize)
	scanner.Buffer(buf, maxBufSize)

	for scanner.Scan() {
		words = append(words, strings.TrimSpace(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}
