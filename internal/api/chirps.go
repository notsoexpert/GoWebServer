package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/notsoexpert/gowebserver/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body,omitempty"`
	UserID    uuid.UUID `json:"user_id"`
	Error     string    `json:"error,omitempty"`
}

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

	sqlChirp, err := cfg.DBQueries.PostChirp(request.Context(), database.PostChirpParams{
		Body:   cleanResponseBody(params.Body),
		UserID: uuid.NullUUID{UUID: params.UserID, Valid: true},
	})
	if err != nil {
		respondWithError(response, 400, "Server failed to create chirp record")
		return
	}

	respBody := Chirp{
		ID:        sqlChirp.ID,
		CreatedAt: sqlChirp.CreatedAt,
		UpdatedAt: sqlChirp.UpdatedAt,
		Body:      sqlChirp.Body,
		UserID:    sqlChirp.UserID.UUID,
	}

	data, encErr := json.Marshal(respBody)
	if encErr != nil {
		respondWithError(response, 500, "Server failed to encode response")
		return
	}
	respondWithJSON(response, 201, data)
}
