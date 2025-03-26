package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	token, err := MakeJWT(uuid.New(), "secret", time.Hour)
	if err != nil {
		t.Errorf("MakeJWT() error = %v", err)
		return
	}
	if token == "" {
		t.Errorf("MakeJWT() token = %v", token)
		return
	}
	t.Logf("MakeJWT() = %v", token)
}

func TestValidateJWT(t *testing.T) {
	tokenUserID := uuid.New()
	secretToken := "secret"
	token, err := MakeJWT(tokenUserID, secretToken, time.Hour)
	if err != nil {
		t.Errorf("MakeJWT() error = %v", err)
		return
	}
	if token == "" {
		t.Errorf("MakeJWT() token = %v", token)
		return
	}
	t.Logf("MakeJWT() = %v", token)

	userID, err := ValidateJWT(token, secretToken)
	if err != nil {
		t.Errorf("ValidateJWT() error = %v", err)
		return
	}
	if userID != tokenUserID {
		t.Errorf("ValidateJWT() userID = %v, want %v", userID, tokenUserID)
		return
	}
}
