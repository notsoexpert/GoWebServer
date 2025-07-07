package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/notsoexpert/gowebserver/internal/auth"
	"github.com/notsoexpert/gowebserver/internal/database"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
		ID:          sqlUser.ID,
		CreatedAt:   sqlUser.CreatedAt,
		UpdatedAt:   sqlUser.UpdatedAt,
		Email:       sqlUser.Email,
		IsChirpyRed: sqlUser.IsChirpyRed,
	}
	data, encErr := json.Marshal(user)
	if encErr != nil {
		respondWithError(response, 500, "Server failed to encode response")
		return
	}
	respondWithJSON(response, 201, data)
}

func (cfg *APIConfig) UpdateCredentialsHandler(response http.ResponseWriter, request *http.Request) {
	accessToken, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(response, 401, "Malformed request")
		return
	}

	validatedID, err := auth.ValidateJWT(accessToken, cfg.Secret)
	if err != nil {
		respondWithError(response, 401, fmt.Sprintf("Unauthorized - %v", err.Error()))
		return
	}

	decoder := json.NewDecoder(request.Body)
	params := credentials{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(response, 401, "Malformed request")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(response, 400, "Could not hash password")
		return
	}

	err = cfg.DBQueries.UpdateUserEmail(request.Context(), database.UpdateUserEmailParams{
		ID:    validatedID,
		Email: params.Email,
	})
	if err != nil {
		respondWithError(response, 500, "Server failed to update credentials")
		return
	}

	err = cfg.DBQueries.UpdateUserPassword(request.Context(), database.UpdateUserPasswordParams{
		ID:             validatedID,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(response, 500, "Server failed to update credentials")
		return
	}

	sqlUser, err := cfg.DBQueries.GetUser(request.Context(), validatedID)
	if err != nil {
		respondWithError(response, 400, "Server failed to retrieve user")
		return
	}

	var user User = User{
		ID:          sqlUser.ID,
		CreatedAt:   sqlUser.CreatedAt,
		UpdatedAt:   sqlUser.UpdatedAt,
		Email:       sqlUser.Email,
		IsChirpyRed: sqlUser.IsChirpyRed,
	}
	data, encErr := json.Marshal(user)
	if encErr != nil {
		respondWithError(response, 500, "Server failed to encode response")
		return
	}
	respondWithJSON(response, 200, data)

}

func (cfg *APIConfig) LoginHandler(response http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	params := credentials{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(response, 400, "Something went wrong")
		return
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

	token, err := auth.MakeJWT(sqlUser.ID, cfg.Secret, 1*time.Hour)
	if err != nil {
		respondWithError(response, 500, "Server failed to authorize token")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(response, 500, "Server failed to authorize token")
		return
	}
	sqlRefreshToken, err := cfg.DBQueries.CreateRefreshToken(request.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    uuid.NullUUID{UUID: sqlUser.ID, Valid: true},
		ExpiresAt: time.Now().Add(60 * time.Hour * 24),
	})
	if err != nil {
		respondWithError(response, 500, "Server failed to authorize token")
		return
	}

	var user User = User{
		ID:           sqlUser.ID,
		CreatedAt:    sqlUser.CreatedAt,
		UpdatedAt:    sqlUser.UpdatedAt,
		Email:        sqlUser.Email,
		Token:        token,
		RefreshToken: sqlRefreshToken.Token,
		IsChirpyRed:  sqlUser.IsChirpyRed,
	}
	data, encErr := json.Marshal(user)
	if encErr != nil {
		respondWithError(response, 500, "Server failed to encode response")
		return
	}
	respondWithJSON(response, 200, data)
}

func (cfg *APIConfig) RefreshHandler(response http.ResponseWriter, request *http.Request) {
	refreshToken, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(response, 401, "Malformed request")
		return
	}

	sqlRefreshToken, err := cfg.DBQueries.GetRefreshToken(request.Context(), refreshToken)
	if err != nil {
		respondWithError(response, 401, "Token not found")
		return
	}
	if time.Now().After(sqlRefreshToken.ExpiresAt) {
		respondWithError(response, 401, "Token expired")
		return
	}
	if sqlRefreshToken.RevokedAt.Valid {
		respondWithError(response, 401, "Token revoked")
	}

	newAccessToken, err := auth.MakeJWT(sqlRefreshToken.UserID.UUID, cfg.Secret, 1*time.Hour)
	if err != nil {
		respondWithError(response, 500, "Server failed to authorize token")
		return
	}

	type AccessTokenResponse struct {
		Token string `json:"token"`
	}
	acToken := AccessTokenResponse{
		Token: newAccessToken,
	}
	data, encErr := json.Marshal(acToken)
	if encErr != nil {
		respondWithError(response, 500, "Server failed to encode response")
		return
	}
	respondWithJSON(response, 200, data)
}

func (cfg *APIConfig) RevokeHandler(response http.ResponseWriter, request *http.Request) {
	refreshToken, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(response, 401, "Malformed request")
		return
	}

	_, err = cfg.DBQueries.GetRefreshToken(request.Context(), refreshToken)
	if err != nil {
		respondWithError(response, 401, "Token not found")
		return
	}

	err = cfg.DBQueries.RevokeRefreshToken(request.Context(), refreshToken)
	if err != nil {
		respondWithError(response, 500, "Server failed to revoke token")
		return
	}
	response.WriteHeader(204)
}
