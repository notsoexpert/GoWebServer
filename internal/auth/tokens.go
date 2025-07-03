package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{0}, err
	}
	issuer, err := token.Claims.(*jwt.RegisteredClaims).GetIssuer()
	if err != nil {
		return uuid.UUID{0}, err
	}
	if issuer != "chirpy" {
		return uuid.UUID{0}, errors.New("invalid token")
	}
	subject, err := token.Claims.(*jwt.RegisteredClaims).GetSubject()
	if err != nil {
		return uuid.UUID{0}, err
	}
	id, err := uuid.Parse(subject)
	if err != nil {
		return uuid.UUID{0}, err
	}
	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authStr := headers.Get("Authorization")
	if authStr == "" {
		return "", errors.New("auth header missing")
	}
	authStr, ok := strings.CutPrefix(authStr, "Bearer ")
	if !ok {
		return "", errors.New("auth header malformed")
	}
	return authStr, nil
}

func MakeRefreshToken() (string, error) {
	return "", nil
}
