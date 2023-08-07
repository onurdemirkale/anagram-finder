package api

import (
	"encoding/json"
	"net/http"
)

type AnagramResponse struct {
	AnagramGroups [][]string `json:"anagramGroups"`
	Error         string     `json:"error,omitempty"`
}

type ResponseServer struct{}

func serveResponse(w http.ResponseWriter, anagramGroups [][]string, status int, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := AnagramResponse{
		AnagramGroups: anagramGroups,
		Error:         err,
	}

	json.NewEncoder(w).Encode(resp)
}
