package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	SALT_ROUNDS := 14
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), SALT_ROUNDS)
	if err != nil {
		return "", err
	}
	hash := string(encryptedPassword)
	return hash, nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
