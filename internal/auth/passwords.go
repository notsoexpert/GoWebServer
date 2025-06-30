package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

/*
Hash the password using the bcrypt.GenerateFromPassword function.
Bcrypt is a secure hash function that is intended for use with passwords.
*/
func HashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", errors.New("attempting to hash password of 0 length")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

/*
Use the bcrypt.CompareHashAndPassword function to compare the password that
the user entered in the HTTP request with the password that is stored in the database.
*/
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
