package main

import (
	"net/http"

	"github.com/onurdemirkale/anagram-finder/api"
	"github.com/onurdemirkale/anagram-finder/pkg/anagram"
	"github.com/onurdemirkale/anagram-finder/pkg/inputsource"
)

func main() {
	var isf inputsource.InputSourceFactoryInterface = inputsource.NewInputSourceFactory()
	aff := anagram.NewAnagramFinderFactory()
	handler := api.NewAnagramHandler(isf, aff)

	http.HandleFunc("/anagram", handler.FindAnagrams)
	http.ListenAndServe(":8080", nil)

}
