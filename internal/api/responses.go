package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

type requestParameters struct {
	Body string `json:"body"`
}

type responseParameters struct {
	CleanedBody string `json:"cleaned_body,omitempty"`
	Error       string `json:"error,omitempty"`
	Valid       bool   `json:"valid"`
}

func cleanResponseBody(body string) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		switch strings.ToLower(word) {
		case "kerfuffle", "sharbert", "fornax":
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func respondWithError(response http.ResponseWriter, code int, msg string) {
	respBody := responseParameters{
		Error: msg,
		Valid: false,
	}
	data, encErr := json.Marshal(respBody)
	if encErr != nil {
		respondWithError(response, 500, "Server failed to encode response") // funny infinite error recursion
		return
	}
	response.WriteHeader(code)
	response.Write(data)
}

func respondWithJSON(response http.ResponseWriter, code int, payload interface{}) {
	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(code)
	response.Write(payload.([]byte))
}
