package utils

import (
	"testing"
)

var newToken string

func TestGenerateToken(t *testing.T) {
	email := "test@gmail.com"
	var id uint = 1
	userId := &id

	token, err := GenerateToken(email, userId)
	newToken = token
	if err != nil {
		t.Errorf("Wrong token")
	}
}

func TestVerifyToken(t *testing.T) {
	_, err := VerifyToken(newToken)
	if err != nil {
		t.Errorf("Wrong token")
	}
}
