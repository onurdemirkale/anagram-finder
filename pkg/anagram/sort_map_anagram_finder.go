package anagram

import (
	"sort"
	"strings"
)

// SortMapAnagramFinder implements the AnagramFinder interface using a basic sort & map approach.
type SortMapAnagramFinder struct{}

func NewSortMapAnagramFinder() *SortMapAnagramFinder {
	return &SortMapAnagramFinder{}
}

// Finds anagrams among the words provided.
// Returns a 2D slice where each sub-slice is a group of anagrams.
// Time complexity: O(N*M*log(M)) where N is the number of words and M is the maximum length of a word.
// Space complexity: O(N*M), the size of the output structure.
func (b *SortMapAnagramFinder) FindAnagrams(words []string) ([][]string, error) {
	anagramGroups := make(map[string][]string)

	for _, word := range words {

		sortedWord := sortWord(word)
		anagramGroups[sortedWord] = append(anagramGroups[sortedWord], word)
	}

	result := make([][]string, 0, len(anagramGroups))

	for _, group := range anagramGroups {
		if len(group) > 1 {
			result = append(result, group)
		}
	}

	return result, nil
}

// todo: handles both normalization and sorting, split into two separate functions
// sortWord takes a word as input, removes spaces, converts it to lowercase, and returns the sorted string.
// This function is used as a helper to normalize the words for anagram comparison.
// Time complexity: O(M*log(M)).
// Space complexity: O(M), where M is the length of the word.
func sortWord(word string) string {
	word = strings.ReplaceAll(word, " ", "")
	word = strings.ToLower(word)

	letters := strings.Split(word, "")
	sort.Strings(letters)

	return strings.Join(letters, "")
}
