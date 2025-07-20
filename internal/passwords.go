package internal

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GenerateTheHashPassword(password string) (string, error) {
	passwordByte := []byte(password) // bcrypt works with only bytes and also returns the hash in bytes.
	bytes, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("error while generating the hash password...", err)
		return "", fmt.Errorf("error")
	}
	return string(bytes), nil
}

func ValidateThePassword(encrypted string, password string) bool {
	passwordCheckErr := bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(password))
	return passwordCheckErr == nil
}
