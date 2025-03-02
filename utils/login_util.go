package utils

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func IsPasswordCorrect(hashed string, pass string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pass))
	if err != nil {
		return false
	} else {
		return true
	}
}

func IsValidID(id string) bool {
	allowedChars := `^[a-zA-Z0-9]`
	regex := regexp.MustCompile(allowedChars)
	return regex.MatchString(id)
}
func IsValidPassword(password string) bool {
	allowedChars := `^[!@#$%^&*()_+=a-zA-Z0-9-]+$`
	regex := regexp.MustCompile(allowedChars)
	return regex.MatchString(password)
}
