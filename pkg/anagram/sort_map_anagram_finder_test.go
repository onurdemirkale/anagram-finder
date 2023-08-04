package anagram

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortMapAnagramFinder_FindAnagrams(t *testing.T) {
	testCases := []struct {
		name     string
		words    []string
		expected [][]string
	}{
		{
			name:     "no words",
			words:    []string{},
			expected: [][]string{},
		},
		{
			name:     "no anagrams",
			words:    []string{"hello", "world"},
			expected: [][]string{},
		},
		{
			name:     "multiple anagrams",
			words:    []string{"cat", "dog", "tac", "god", "good", "act"},
			expected: [][]string{{"cat", "tac", "act"}, {"dog", "god"}},
		},
		{
			name:     "punctuations or symbols are not considered",
			words:    []string{"hello!", "world??"},
			expected: [][]string{},
		},
		{
			name:     "case sensitivity",
			words:    []string{"Cat", "tac"},
			expected: [][]string{{"Cat", "tac"}},
		},
		{
			name:     "multi-word anagrams",
			words:    []string{"debit card", "bad credit", "cat", "tac"},
			expected: [][]string{{"debit card", "bad credit"}, {"cat", "tac"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			smaf := NewSortMapAnagramFinder()
			actual, _ := smaf.FindAnagrams(tc.words)

			// sort the anagram groups to ensure ElementsMatch works correctly
			for i := range actual {
				sort.Strings(actual[i])
			}

			for i := range tc.expected {
				sort.Strings(tc.expected[i])
			}

			assert.ElementsMatch(t, tc.expected, actual)
		})
	}
}
