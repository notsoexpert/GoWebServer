package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/notsoexpert/gowebserver/internal/auth"
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

func ReadyChirpForJSON(sqlChirp database.Chirp) Chirp {
	return Chirp{
		ID:        sqlChirp.ID,
		CreatedAt: sqlChirp.CreatedAt,
		UpdatedAt: sqlChirp.UpdatedAt,
		Body:      sqlChirp.Body,
		UserID:    sqlChirp.UserID.UUID,
	}
}

func (cfg *APIConfig) PostChirpsHandler(response http.ResponseWriter, request *http.Request) {
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

	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(response, 401, "Malformed request")
		return
	}

	validatedID, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		respondWithError(response, 401, fmt.Sprintf("Unauthorized - %v", err.Error()))
		return
	}

	if len(params.Body) > 140 {
		respondWithError(response, 400, "Chirp is too long")
		return
	}

	sqlChirp, err := cfg.DBQueries.PostChirp(request.Context(), database.PostChirpParams{
		Body:   cleanResponseBody(params.Body),
		UserID: uuid.NullUUID{UUID: validatedID, Valid: true},
	})
	if err != nil {
		respondWithError(response, 400, "Server failed to create chirp record")
		return
	}

	respBody := ReadyChirpForJSON(sqlChirp)

	data, encErr := json.Marshal(respBody)
	if encErr != nil {
		respondWithError(response, 500, "Server failed to encode response")
		return
	}
	respondWithJSON(response, 201, data)
}

func (cfg *APIConfig) GetChirpsHandler(response http.ResponseWriter, request *http.Request) {
	sqlChirps, err := cfg.DBQueries.GetChirps(request.Context())
	if err != nil {
		respondWithError(response, 400, "Server failed to get chirp records")
		return
	}

	var respBody []Chirp

	for _, sqlChirp := range sqlChirps {
		respBody = append(respBody, ReadyChirpForJSON(sqlChirp))
	}

	data, encErr := json.Marshal(respBody)
	if encErr != nil {
		respondWithError(response, 500, "Server failed to encode response")
		return
	}
	respondWithJSON(response, 200, data)
}

func (cfg *APIConfig) GetChirpHandler(response http.ResponseWriter, request *http.Request) {
	sqlChirps, err := cfg.DBQueries.GetChirps(request.Context())
	if err != nil {
		respondWithError(response, 400, "Server failed to get chirp records")
		return
	}

	uuid, err := uuid.Parse(request.PathValue("chirpID"))
	if err != nil {
		respondWithError(response, 404, "Chirp not found")
		return
	}

	var foundChirp Chirp
	err = errors.New("Chirp not found")
	for _, sqlChirp := range sqlChirps {
		chirp := ReadyChirpForJSON(sqlChirp)

		if chirp.ID == uuid {
			foundChirp = chirp
			err = nil
			break
		}
	}

	if err != nil {
		respondWithError(response, 404, err.Error())
		return
	}

	data, encErr := json.Marshal(foundChirp)
	if encErr != nil {
		respondWithError(response, 500, "Server failed to encode response")
		return
	}
	respondWithJSON(response, 200, data)
}
