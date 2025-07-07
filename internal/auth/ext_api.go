package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authStr := headers.Get("Authorization")
	if authStr == "" {
		return "", errors.New("auth header missing")
	}
	authStr, ok := strings.CutPrefix(authStr, "ApiKey ")
	if !ok {
		return "", errors.New("auth header malformed")
	}
	return authStr, nil
}
