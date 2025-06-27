package api

import (
	"encoding/json"
	"net/http"
)

func ValidateChirpHandler(response http.ResponseWriter, request *http.Request) {
	type requestParameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(request.Body)
	params := requestParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(response, 400, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(response, 400, "Chirp is too long")
		return
	}

	respBody := responseParameters{
		CleanedBody: cleanResponseBody(params.Body),
		Valid:       true,
	}
	data, encErr := json.Marshal(respBody)
	if encErr != nil {
		respondWithError(response, 500, "Server failed to encode response")
		return
	}
	respondWithJSON(response, 200, data)
}
