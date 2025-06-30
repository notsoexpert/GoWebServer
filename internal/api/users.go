package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *APIConfig) CreateUserHandler(response http.ResponseWriter, request *http.Request) {
	type requestParameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(request.Body)
	params := requestParameters{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(response, 400, "Something went wrong")
		return
	}

	sqlUser, err := cfg.DBQueries.CreateUser(request.Context(), params.Email)
	if err != nil {
		respondWithError(response, 400, "Server failed to create user")
		return
	}

	var user User = User{
		ID:        sqlUser.ID,
		CreatedAt: sqlUser.CreatedAt,
		UpdatedAt: sqlUser.UpdatedAt,
		Email:     sqlUser.Email,
	}
	data, encErr := json.Marshal(user)
	if encErr != nil {
		respondWithError(response, 500, "Server failed to encode response")
		return
	}
	respondWithJSON(response, 201, data)
}
