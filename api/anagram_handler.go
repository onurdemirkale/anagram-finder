package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/onurdemirkale/anagram-finder/pkg/anagram"
	"github.com/onurdemirkale/anagram-finder/pkg/inputsource"
)

const (
	errProcessing             = "an error occurred while processing your request"
	errInvalidFormat          = "invalid request format"
	errInvalidFile            = "failed to read file"
	errUnsupportedContentType = "unsupported content type"
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

func logAndRespondError(w http.ResponseWriter, err error, status int, userErrorMessage string) {
	log.Printf("handler error: %v", err)
	http.Error(w, userErrorMessage, status)
}

func (h *AnagramHandler) FindAnagrams(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	inputSource, req, err := h.parseRequest(r)
	if err != nil {
		logAndRespondError(w, err, http.StatusBadRequest, err.Error())
		return
	}

	err = req.Validate()
	if err != nil {
		logAndRespondError(w, err, http.StatusBadRequest, err.Error())
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

	if req.InputType == inputTypeFile {
		h.serveAnagramsAsFile(w, r, anagramGroups)
	} else {
		h.serveAnagramsAsJSON(w, anagramGroups)
	}
}

func (h *AnagramHandler) parseRequest(r *http.Request) (inputsource.InputSource, AnagramRequest, error) {
	var req AnagramRequest
	contentType := r.Header.Get("Content-Type")

	switch {
	case strings.Contains(contentType, "multipart/form-data"):
		file, _, err := r.FormFile("file")
		if err != nil {
			return nil, req, fmt.Errorf("%s: %v", errInvalidFile, err)
		}
		defer file.Close()

		req.InputType = r.FormValue("inputType")
		req.Algorithm = r.FormValue("algorithm")

		return inputsource.NewHttpFileInputSource(file), req, nil

	case strings.Contains(contentType, "application/json"):
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			return nil, req, fmt.Errorf("%s: %v", errInvalidFormat, err)
		}

		inputSource, err := h.inputSourceFactory.CreateInputSource(req.InputType, req.InputData)
		if err != nil {
			return nil, req, err
		}

		return inputSource, req, nil

	default:
		return nil, req, fmt.Errorf(errUnsupportedContentType)
	}
}

func (h *AnagramHandler) serveAnagramsAsFile(w http.ResponseWriter, r *http.Request, groups [][]string) {
	tmpFile, err := os.CreateTemp("", "anagram-output-*.txt")
	if err != nil {
		logAndRespondError(w, err, http.StatusInternalServerError, errProcessing)
		return
	}
	defer os.Remove(tmpFile.Name())

	for _, group := range groups {
		line := strings.Join(group, ", ")
		fmt.Fprintln(tmpFile, line)
	}
	tmpFile.Close()

	w.Header().Set("Content-Disposition", "attachment; filename=anagram-output.txt")
	w.Header().Set("Content-Type", "text/plain")
	http.ServeFile(w, r, tmpFile.Name())
}

func (h *AnagramHandler) serveAnagramsAsJSON(w http.ResponseWriter, groups [][]string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}
