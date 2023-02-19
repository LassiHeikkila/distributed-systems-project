package main

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pw string) string {
	const cost = 12
	b, err := bcrypt.GenerateFromPassword([]byte(pw), cost)
	if err != nil {
		return ""
	}
	return string(b)
}

func PasswordEqualsHash(plain string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil
}
