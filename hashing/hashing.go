package hashing

import (
	"golang.org/x/crypto/bcrypt"
)

func CreateHash(passwordString string) (hashPassword string, errInHashPsw error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwordString), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func VarifyPassword(hashPassword, passwordString string) (IsMatched bool) {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(passwordString))
	return err == nil
}
