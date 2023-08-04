package anagram

type AnagramFinder interface {
	FindAnagrams(words []string) ([][]string, error)
}
