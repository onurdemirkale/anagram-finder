package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/onurdemirkale/anagram-finder/pkg/anagram"
	"github.com/onurdemirkale/anagram-finder/pkg/inputsource"
)

const (
	errProcessing    = "an error occurred while processing your request"
	errInvalidFormat = "invalid request format"
)

type AnagramHandler struct {
	inputSourceFactory   inputsource.InputSourceFactoryInterface
	anagramFinderFactory anagram.AnagramFinderFactoryInterface
}

func NewAnagramHandler(isf inputsource.InputSourceFactoryInterface, aff anagram.AnagramFinderFactoryInterface) *AnagramHandler {
	return &AnagramHandler{
		inputSourceFactory:   isf,
		anagramFinderFactory: aff,
	}
}

func decodeAndValidateRequest(r *http.Request) (AnagramRequest, error) {
	var req AnagramRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, fmt.Errorf("%s: %v", errInvalidFormat, err)
	}

	if err := req.Validate(); err != nil {
		return req, err
	}

	return req, nil
}

func logAndRespondError(w http.ResponseWriter, err error, status int, userErrorMessage string) {
	log.Printf("handler error: %v", err)
	http.Error(w, userErrorMessage, status)
}

func (h *AnagramHandler) FindAnagrams(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	req, err := decodeAndValidateRequest(r)
	if err != nil {
		logAndRespondError(w, err, http.StatusBadRequest, err.Error())
		return
	}

	inputSource, err := h.inputSourceFactory.CreateInputSource(req.InputType, req.InputData)
	if err != nil {
		logAndRespondError(w, err, http.StatusBadRequest, errProcessing)
		return
	}

	anagramFinder, err := h.anagramFinderFactory.CreateAnagramFinder(req.Algorithm)
	if err != nil {
		logAndRespondError(w, err, http.StatusBadRequest, errProcessing)
		return
	}

	words, err := inputSource.GetWords()
	if err != nil {
		logAndRespondError(w, err, http.StatusInternalServerError, errProcessing)
		return
	}

	anagramGroups, err := anagramFinder.FindAnagrams(words)
	if err != nil {
		logAndRespondError(w, err, http.StatusInternalServerError, errProcessing)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(anagramGroups); err != nil {
		logAndRespondError(w, err, http.StatusInternalServerError, errProcessing)
		return
	}
}
