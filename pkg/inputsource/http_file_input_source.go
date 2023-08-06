package inputsource

import (
	"bufio"
	"mime/multipart"
)

type HttpFileInputSource struct {
	file multipart.File
}

func NewHttpFileInputSource(file multipart.File) *HttpFileInputSource {
	return &HttpFileInputSource{file: file}
}

func (hf *HttpFileInputSource) GetWords() ([]string, error) {
	defer hf.file.Close()

	var words []string

	scanner := bufio.NewScanner(hf.file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}
