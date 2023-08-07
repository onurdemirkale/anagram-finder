package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/onurdemirkale/anagram-finder/pkg/anagram"
	"github.com/onurdemirkale/anagram-finder/pkg/inputsource"
)

func TestFindAnagrams(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		expectedCode   int
		expectedOutput string
		expectedError  string
	}{
		{
			name:           "Valid Request",
			body:           `{"inputType": "http", "inputData": "listen,enlist,inlets,cat,silent,tac,nag a ram,anagram", "algorithm": "basic"}`,
			expectedCode:   http.StatusOK,
			expectedOutput: "{\"anagramGroups\":[[\"listen\",\"enlist\",\"inlets\",\"silent\"],[\"cat\",\"tac\"],[\"nag a ram\",\"anagram\"]]}\n",
		},
		{
			name:          "Invalid Input Type",
			body:          `{"inputType": "invalid", "inputData": "listen,enlist,inlets,silent", "algorithm": "basic"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: fmt.Sprintf("{\"anagramGroups\":null,\"error\":\"%s\"}", ErrInvalidInputType),
		},
		{
			name:          "Invalid Algorithm",
			body:          `{"inputType": "http", "inputData": "listen,enlist,inlets,silent", "algorithm": "invalid"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: fmt.Sprintf("{\"anagramGroups\":null,\"error\":\"%s\"}", ErrInvalidAlgorithmType),
		},
		{
			name:          "Invalid Input Data",
			body:          `{"inputType": "http", "inputData": "invalidData", "algorithm": "basic"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: fmt.Sprintf("{\"anagramGroups\":null,\"error\":\"%s\"}", ErrInvalidInput),
		},
	}

	handler := NewAnagramHandler(&inputsource.InputSourceFactory{}, &anagram.AnagramFinderFactory{})

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/anagram", bytes.NewBuffer([]byte(tc.body)))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.FindAnagrams(rr, req)

			if status := rr.Code; status != tc.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedCode)
			}

			actual := rr.Body.String()

			if tc.expectedError != "" {
				actualError := strings.TrimSpace(actual)

				if actualError != tc.expectedError {
					t.Errorf("handler returned unexpected error: got %q want %q", actualError, tc.expectedError)
				}
			} else {
				actualSorted := sortJsonResponse(actual)
				expectedSorted := sortJsonResponse(tc.expectedOutput)
				if actualSorted != expectedSorted {
					t.Errorf("handler returned unexpected file content: got %q want %q", actualSorted, expectedSorted)
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
		expectedError      error
		expectedFileOutput string
	}{
		{
			name:               "Valid File Input",
			fileContents:       "listen\nenlist\ninlets\ncat\nsilent\ntac\nnag a ram\nanagram",
			inputType:          "http_file",
			algorithm:          "basic",
			expectedCode:       http.StatusOK,
			expectedError:      nil,
			expectedFileOutput: "{\"anagramGroups\":[[\"listen\",\"enlist\",\"inlets\",\"silent\"],[\"cat\",\"tac\"],[\"nag a ram\",\"anagram\"]]}\n",
		},
	}

	handler := NewAnagramHandler(&inputsource.InputSourceFactory{}, &anagram.AnagramFinderFactory{})

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := generateMultipartRequest(tc.fileContents, tc.inputType, tc.algorithm)
			rr := httptest.NewRecorder()

			handler.FindAnagrams(rr, req)

			if status := rr.Code; status != tc.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedCode)
			}

			actual := rr.Body.String()

			if tc.expectedError != nil {
				actualErr := strings.TrimSpace(actual)
				expectedErr := tc.expectedError.Error()

				if actualErr != expectedErr {
					t.Errorf("handler returned unexpected error: got %q want %q", actualErr, expectedErr)
				}
			}

			actualSorted := sortJsonResponse(actual)
			expectedSorted := sortJsonResponse(tc.expectedFileOutput)
			if actualSorted != expectedSorted {
				t.Errorf("handler returned unexpected file content: got %q want %q", actualSorted, expectedSorted)
			}
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

func sortJsonResponse(jsonStr string) string {
	var response AnagramResponse
	json.Unmarshal([]byte(jsonStr), &response)

	// sort the anagrams within each group
	for _, group := range response.AnagramGroups {
		sort.Strings(group)
	}

	// sort the groups themselves
	sort.Slice(response.AnagramGroups, func(i, j int) bool {
		return strings.Join(response.AnagramGroups[i], ",") < strings.Join(response.AnagramGroups[j], ",")
	})

	output, _ := json.Marshal(response)
	return string(output)
}
