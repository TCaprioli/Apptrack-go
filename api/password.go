package api

import (
	"golang.org/x/crypto/bcrypt"
)

var MinimumPasswordLength = 8
var MaximumPasswordLength = 32

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
