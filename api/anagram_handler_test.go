package api

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/onurdemirkale/anagram-finder/pkg/anagram"
	"github.com/onurdemirkale/anagram-finder/pkg/inputsource"
)

// todo: parameterize mocks
type MockInputSource struct{}

func (m *MockInputSource) GetWords() ([]string, error) {
	return []string{"listen", "enlist", "inlets", "cat", "silent", "tac", "nag a ram", "anagram"}, nil
}

type MockInputSourceFactory struct{}

func (m *MockInputSourceFactory) CreateInputSource(inputType string, inputData string) (inputsource.InputSource, error) {
	if inputType == "error" {
		return nil, errors.New("InputSourceFactory error")
	}
	return &MockInputSource{}, nil
}

type MockAnagramFinder struct{}

func (m *MockAnagramFinder) FindAnagrams(words []string) ([][]string, error) {
	return [][]string{{"listen", "enlist", "inlets", "silent"}, {"cat", "tac"}, {"nag a ram", "anagram"}}, nil
}

type MockAnagramFinderFactory struct{}

func (m *MockAnagramFinderFactory) CreateAnagramFinder(algo string) (anagram.AnagramFinder, error) {
	if algo == "error" {
		return nil, errors.New("AnagramFinderFactory error")
	}
	return &MockAnagramFinder{}, nil
}

func TestFindAnagrams(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		expectedCode   int
		expectedOutput string
		expectedError  error
	}{
		{
			name:           "Valid Request",
			body:           `{"inputType": "http", "inputData": "listen,enlist,inlets,cat,silent,tac", "algorithm": "basic"}`,
			expectedCode:   http.StatusOK,
			expectedOutput: `[["listen","enlist","inlets","silent"],["cat","tac"],["nag a ram","anagram"]]`,
		},
		{
			name:          "Invalid Input Type",
			body:          `{"inputType": "invalid", "inputData": "listen,enlist,inlets,silent", "algorithm": "basic"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: errInvalidInputType,
		},
		{
			name:          "Invalid Algorithm",
			body:          `{"inputType": "http", "inputData": "listen,enlist,inlets,silent", "algorithm": "invalid"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: errInvalidAlgorithmType,
		},
		{
			name:          "Invalid Input Data",
			body:          `{"inputType": "http", "inputData": "invalidData", "algorithm": "basic"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: errInvalidInput,
		},
	}

	handler := NewAnagramHandler(&MockInputSourceFactory{}, &MockAnagramFinderFactory{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/anagram", bytes.NewBuffer([]byte(tt.body)))
			rr := httptest.NewRecorder()

			handler.FindAnagrams(rr, req)

			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedCode)
			}

			actual := strings.TrimSpace(rr.Body.String())

			if tt.expectedError != nil {
				expectedError := tt.expectedError.Error()

				if actual != expectedError {
					t.Errorf("handler returned unexpected error: got %q want %q", actual, expectedError)
				}
			} else if tt.expectedOutput != "" {
				if actual != tt.expectedOutput {
					t.Errorf("handler returned unexpected body: got %q want %q", actual, tt.expectedOutput)
				}
			}
		})
	}
}
