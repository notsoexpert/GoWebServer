package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/notsoexpert/gowebserver/internal/auth"
	"github.com/notsoexpert/gowebserver/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

type credentials struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	ExpiresIn int64  `json:"expires_in_seconds,omitempty"`
}

func (cfg *APIConfig) CreateUserHandler(response http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	params := credentials{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(response, 400, "Something went wrong")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(response, 400, "Could not hash password")
		return
	}

	sqlUser, err := cfg.DBQueries.CreateUser(request.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
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

func (cfg *APIConfig) LoginHandler(response http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	params := credentials{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(response, 400, "Something went wrong")
		return
	}

	expiresIn := 1 * time.Hour
	if params.ExpiresIn != 0 && (time.Duration(params.ExpiresIn)*time.Second) > time.Hour {
		expiresIn = time.Duration(params.ExpiresIn) * time.Second
	}

	sqlUser, err := cfg.DBQueries.GetUserByEmail(request.Context(), params.Email)
	if err != nil {
		respondWithError(response, 401, "Incorrect email or password")
		return
	}

	if err := auth.CheckPasswordHash(params.Password, sqlUser.HashedPassword); err != nil {
		respondWithError(response, 401, "Incorrect email or password")
		return
	}

	token, err := auth.MakeJWT(sqlUser.ID, cfg.Secret, expiresIn)
	if err != nil {
		respondWithError(response, 500, "Server failed to authorize token")
	}

	var user User = User{
		ID:        sqlUser.ID,
		CreatedAt: sqlUser.CreatedAt,
		UpdatedAt: sqlUser.UpdatedAt,
		Email:     sqlUser.Email,
		Token:     token,
	}
	data, encErr := json.Marshal(user)
	if encErr != nil {
		respondWithError(response, 500, "Server failed to encode response")
		return
	}
	respondWithJSON(response, 200, data)
}
