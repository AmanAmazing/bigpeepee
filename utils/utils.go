package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(plainPass string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPass), 14)
	return string(hashedPassword), err
}

func CheckPasswordMatch(hashedPassword, enteredPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(enteredPassword))
	return err == nil
}
