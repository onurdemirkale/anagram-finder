package anagram

import "errors"

type AnagramFinderFactoryInterface interface {
	CreateAnagramFinder(algorithm string) (AnagramFinder, error)
}

type AnagramFinderFactory struct{}

func NewAnagramFinderFactory() AnagramFinderFactoryInterface {
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
