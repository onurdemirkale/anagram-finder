package inputsource

import (
	"bytes"
	"reflect"
	"testing"
)

type MockMultipartFile struct {
	content *bytes.Reader
}

func (m *MockMultipartFile) Read(p []byte) (n int, err error) {
	return m.content.Read(p)
}

func (m *MockMultipartFile) ReadAt(p []byte, off int64) (n int, err error) {
	return m.content.ReadAt(p, off)
}

func (m *MockMultipartFile) Seek(offset int64, whence int) (int64, error) {
	return m.content.Seek(offset, whence)
}

func (m *MockMultipartFile) Close() error {
	return nil
}

func NewMockMultipartFile(data string) *MockMultipartFile {
	return &MockMultipartFile{
		content: bytes.NewReader([]byte(data)),
	}
}

func TestHttpFileInputSource_GetWords(t *testing.T) {
	tests := []struct {
		name         string
		fileContents string
		expected     []string
	}{
		{
			name:         "Multiple lines",
			fileContents: "hello\nworld",
			expected:     []string{"hello", "world"},
		},
		{
			name:         "Single line",
			fileContents: "singleword",
			expected:     []string{"singleword"},
		},
		{
			name:         "Empty file",
			fileContents: "",
			expected:     []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockFile := NewMockMultipartFile(tc.fileContents)
			source := NewHttpFileInputSource(mockFile)
			words, err := source.GetWords()

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if reflect.ValueOf(words).Len() != reflect.ValueOf(tc.expected).Len() {
				t.Errorf("expected words %v, got %v", tc.expected, words)
			}

		})
	}
}
