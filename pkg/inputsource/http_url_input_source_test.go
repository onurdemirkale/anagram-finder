package inputsource

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHttpUrlInputSource_GetWords(t *testing.T) {
	tests := []struct {
		name        string
		serverWords []string
		expected    []string
	}{
		{
			name:        "normal case",
			serverWords: []string{"apple", "banana", "cherry"},
			expected:    []string{"apple", "banana", "cherry"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for _, word := range tc.serverWords {
					fmt.Fprintln(w, word)
				}
			}))
			defer server.Close()

			source := NewHttpUrlInputSource(server.URL)
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
