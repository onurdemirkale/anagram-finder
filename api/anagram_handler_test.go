package api

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
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

func (m *MockInputSourceFactory) CreateInputSource(inputType string, inputData interface{}) (inputsource.InputSource, error) {
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
			req.Header.Set("Content-Type", "application/json")
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

// todo: implement test scenarios for
// large files
// different file types
// malformed content
// empty files
// special characters
func TestFindAnagrams_FileInput(t *testing.T) {
	tests := []struct {
		name               string
		fileContents       string
		inputType          string
		algorithm          string
		expectedCode       int
		expectedErr        error
		expectedFileOutput string
	}{
		{
			name:               "Valid File Input",
			fileContents:       "listen\nenlist\ninlets\ncat\nsilent\ntac\nnag a ram\nanagram",
			inputType:          "http_file",
			algorithm:          "basic",
			expectedCode:       http.StatusOK,
			expectedErr:        nil,
			expectedFileOutput: "listen, enlist, inlets, silent\ncat, tac\nnag a ram, anagram\n",
		},
	}

	handler := NewAnagramHandler(&MockInputSourceFactory{}, &MockAnagramFinderFactory{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := generateMultipartRequest(tt.fileContents, tt.inputType, tt.algorithm)
			rr := httptest.NewRecorder()

			handler.FindAnagrams(rr, req)

			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedCode)
			}

			if tt.expectedErr != nil {
				actualErr := strings.TrimSpace(rr.Body.String())
				expectedErr := tt.expectedErr.Error()

				if actualErr != expectedErr {
					t.Errorf("handler returned unexpected error: got %q want %q", actualErr, expectedErr)
				}
			}

			validateFileContent(t, rr, tt.expectedFileOutput)
		})
	}
}

func generateMultipartRequest(fileContents, inputType, algorithm string) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte(fileContents))

	writer.WriteField("inputType", inputType)
	writer.WriteField("algorithm", algorithm)
	writer.Close()

	req := httptest.NewRequest("POST", "/anagram", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	return req
}

func validateFileContent(t *testing.T, rr *httptest.ResponseRecorder, expectedFileOutput string) {
	tempFile, err := os.CreateTemp("", "test-output-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, rr.Result().Body)
	if err != nil {
		t.Fatal(err)
	}
	rr.Result().Body.Close()

	tempFile.Seek(0, 0)
	fileContentBytes, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if expectedFileOutput != "" {
		actual := string(fileContentBytes)
		if actual != expectedFileOutput {
			t.Errorf("handler returned unexpected file content: got %q want %q", actual, expectedFileOutput)
		}
	}
}
