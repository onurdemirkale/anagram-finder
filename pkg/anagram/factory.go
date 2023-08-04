package anagram

import "errors"

type AnagramFinderFactory struct{}

func NewAnagramFinderFactory() *AnagramFinderFactory {
	return &AnagramFinderFactory{}
}

func (f *AnagramFinderFactory) CreateAnagramFinder(algorithm string) (AnagramFinder, error) {
	switch algorithm {
	case "basic":
		return NewSortMapAnagramFinder(), nil
	default:
		return nil, errors.New("unknown algorithm")
	}
}
