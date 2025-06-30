package auth

import "testing"

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
