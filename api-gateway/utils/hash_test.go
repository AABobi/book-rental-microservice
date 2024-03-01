package utils

import (
	"errors"
	"testing"
)

func TestAbs(t *testing.T) {
	pass, _ := HashPassword("pass")

	if !CheckPasswordHash("pass", pass) {
		t.Fatal(errors.New("CheckPasswordHash error"))
	}
}
