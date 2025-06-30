package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *APIConfig) ChirpsHandler(response http.ResponseWriter, request *http.Request) {
	type requestParameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
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
		Body:   cleanResponseBody(params.Body),
		UserID: params.UserID,
	}
	data, encErr := json.Marshal(respBody)
	if encErr != nil {
		respondWithError(response, 500, "Server failed to encode response")
		return
	}
	respondWithJSON(response, 201, data)
}
