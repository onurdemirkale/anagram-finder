package inputsource

import (
	"reflect"
	"testing"
)

func TestHttpBodyInputSource_GetWords(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "Simple content",
			content:  "word1,word2,word3",
			expected: []string{"word1", "word2", "word3"},
		},
		{
			name:     "Empty content",
			content:  "",
			expected: []string{""},
		},
		{
			name:     "Single word",
			content:  "word1",
			expected: []string{"word1"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			source := NewHttpBodyInputSource(tc.content)
			words, err := source.GetWords()

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(words, tc.expected) {
				t.Errorf("expected words %v, got %v", tc.expected, words)
			}
		})
	}
}
