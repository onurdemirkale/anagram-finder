package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/onurdemirkale/anagram-finder/pkg/anagram"
	"github.com/onurdemirkale/anagram-finder/pkg/inputsource"
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

func (h *AnagramHandler) FindAnagrams(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	inputSource, req, err := h.parseRequest(r)
	if err != nil {
		status, errMsg := handleError(err)
		serveResponse(w, nil, status, errMsg)
		return
	}

	anagramGroups, err := h.processAnagrams(req, inputSource)
	if err != nil {
		status, errMsg := handleError(err)
		serveResponse(w, nil, status, errMsg)
		return
	}

	serveResponse(w, anagramGroups, http.StatusOK, "")
}

// todo: this method does not adhere to SOLID (SRP), refactor
func (h *AnagramHandler) parseRequest(r *http.Request) (inputsource.InputSource, AnagramRequest, error) {
	var req AnagramRequest
	contentType := r.Header.Get("Content-Type")

	switch {
	case strings.Contains(contentType, "multipart/form-data"):
		file, _, err := r.FormFile("file")
		if err != nil {
			return nil, req, errors.New(ErrInvalidFileInput)
		}
		defer file.Close()

		req.InputType = r.FormValue("inputType")
		req.Algorithm = r.FormValue("algorithm")

		if err := req.validate(); err != nil {
			return nil, req, err
		}

		return inputsource.NewHttpFileInputSource(file), req, nil

	case strings.Contains(contentType, "application/json"):
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			return nil, req, errors.New(ErrInvalidFormat)
		}

		if err := req.validate(); err != nil {
			return nil, req, err
		}

		inputSource, err := h.inputSourceFactory.CreateInputSource(req.InputType, req.InputData)
		if err != nil {
			return nil, req, err
		}

		return inputSource, req, nil

	default:
		return nil, req, errors.New(ErrUnsupportedContentType)
	}
}

func (h *AnagramHandler) processAnagrams(req AnagramRequest, inputSource inputsource.InputSource) ([][]string, error) {
	anagramFinder, err := h.anagramFinderFactory.CreateAnagramFinder(req.Algorithm)
	if err != nil {
		return nil, err
	}

	words, err := inputSource.GetWords()
	if err != nil {
		return nil, err
	}

	anagramGroups, err := anagramFinder.FindAnagrams(words)
	if err != nil {
		return nil, err
	}

	return anagramGroups, nil
}
