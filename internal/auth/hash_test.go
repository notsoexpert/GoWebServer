package auth

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "0123456"

	msg, err := HashPassword(password)
	if err != nil {
		t.Errorf(`HashPassword("0123456") failed) = %q, %v`, msg, err)
	}
}

func TestEmptyPassword(t *testing.T) {
	password := ""

	msg, err := HashPassword(password)
	if err == nil {
		t.Errorf(`HashPassword("") succeeded) = %q`, msg)
	}
}

func TestLongPassword(t *testing.T) {
	password := "3.141592653589793238462643383279502884197169399375105820974944592307816406"

	msg, err := HashPassword(password)
	if err == nil {
		t.Errorf(`HashPassword("3.141592653589793238462643383279502884197169399375105820974944592307816406") succeeded) = %q`, msg)
	}
}

func TestHashSamePassword(t *testing.T) {
	password1 := "password"
	password2 := "password"

	hashed1, err := HashPassword(password1)
	if err != nil {
		t.Errorf(`1st HashPassword("password") = %q, %v`, hashed1, err)
	}
	hashed2, err := HashPassword(password2)
	if err != nil {
		t.Errorf(`2nd HashPassword("password") = %q, %v`, hashed2, err)
	}

	if hashed1 == hashed2 {
		t.Errorf(`Same passwords hashed to same value - %q == %q`, hashed1, hashed2)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "password"
	hashed_password, err := HashPassword(password)
	if err != nil {
		t.Errorf(`HashPassword("password") = %q, %v`, hashed_password, err)
	}

	if err := CheckPasswordHash(password, hashed_password); err != nil {
		t.Errorf(`ChechPasswordHash(%q, %q) failed = %v`, password, hashed_password, err)
	}
}

func TestCreateAndValidateJWT(t *testing.T) {
	id := uuid.New()
	tokenSecret := "secret"
	expiresIn := 100 * time.Hour

	tokenString, err := MakeJWT(id, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf(`MakeJWT(%q, %q, %d) failed = %v`, id.String(), tokenSecret, int64(expiresIn), err)
	}

	retID, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Errorf(`ValidateJWT(%q, %q) failed = %v`, tokenString, tokenSecret, err)
	}

	if id != retID {
		t.Errorf(`ValidateJWT(%q, %q) failed - returned uuid %v does not match original %v`, tokenString, tokenSecret, retID.String(), id.String())
	}
}

func TestExpiredJWT(t *testing.T) {
	id := uuid.New()
	tokenSecret := "secret"
	var expiresIn time.Duration

	tokenString, err := MakeJWT(id, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf(`MakeJWT(%q, %q, %d) failed = %v`, id.String(), tokenSecret, int64(expiresIn), err)
	}

	_, err = ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return
		}
		t.Errorf(`ValidateJWT(%q, %q) failed, but not due to expiration = %v`, tokenString, tokenSecret, err)
	}

	t.Errorf(`ValidateJWT(%q, %q) succeeded - expired token passed validation`, tokenString, tokenSecret)
}

func TestWrongSecretJWT(t *testing.T) {
	id := uuid.New()
	tokenSecret := "secret"
	expiresIn := 100 * time.Hour

	tokenString, err := MakeJWT(id, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf(`MakeJWT(%q, %q, %d) failed = %v`, id.String(), tokenSecret, int64(expiresIn), err)
	}

	wrongTokenSecret := "wrongsecret"
	_, err = ValidateJWT(tokenString, wrongTokenSecret)
	if err != nil {
		if strings.Contains(err.Error(), "signature is invalid") {
			return
		}
		t.Errorf(`ValidateJWT(%q, %q) failed, but not due to wrong secret = %v`, tokenString, wrongTokenSecret, err)
	}

	t.Errorf(`ValidateJWT(%q, %q) succeeded - wrong secret token passed validation`, tokenString, tokenSecret)
}

func TestGetBearerToken(t *testing.T) {
	header := make(http.Header)
	header.Set("Authorization", "Bearer QWERTYUIOP")

	tokenString, err := GetBearerToken(header)
	if err != nil {
		t.Errorf(`GetBearerToken failed - %q, %v`, tokenString, err)
	}
	if !strings.Contains(tokenString, "QWERTYUIOP") {
		t.Errorf(`GetBearerToken returned incorrect token string - %q, %v`, tokenString, err)
	}
}
