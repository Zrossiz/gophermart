package utils

import (
	"testing"
	"time"
)

func TestGenerateJWT(t *testing.T) {
	var generateJWTProps GenerateJWTProps = GenerateJWTProps{
		Secret:   []byte("test"),
		Exprires: time.Now(),
		UserID:   1,
		Username: "test",
	}

	_, err := GenerateJWT(generateJWTProps)
	if err != nil {
		t.Fatal("expected error == nil, got: ", err)
	}
}
