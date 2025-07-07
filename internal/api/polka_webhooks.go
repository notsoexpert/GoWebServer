package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/notsoexpert/gowebserver/internal/auth"
)

type parameters struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *APIConfig) PolkaWebhooksHandler(response http.ResponseWriter, request *http.Request) {
	apiKey, err := auth.GetAPIKey(request.Header)
	if err != nil {
		respondWithError(response, 401, "Malformed request")
		return
	}

	if !strings.Contains(cfg.PolkaKey, apiKey) {
		respondWithError(response, 401, "Authentication failure")
		return
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(response, 400, "Malformed request")
		return
	}

	if !strings.Contains(params.Event, "user.upgraded") {
		response.WriteHeader(204)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(response, 400, "Invalid user ID")
		return
	}

	if err := cfg.DBQueries.ActivateChirpyRed(request.Context(), userID); err != nil {
		respondWithError(response, 404, "Invalid user ID")
		return
	}
	response.WriteHeader(204)
}
