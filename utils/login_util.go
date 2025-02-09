package utils

import "golang.org/x/crypto/bcrypt"

func IsPasswordCorrect(hashed string, pass string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pass))
	if err != nil {
		return false
	} else {
		return true
	}
}
