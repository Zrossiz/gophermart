package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type GenerateJWTProps struct {
	secret   []byte
	exprires time.Time
	userID   int64
	username string
}

func GenerateJWT(props GenerateJWTProps) (string, error) {
	claims := &CustomClaims{
		UserID:   props.userID,
		Username: props.username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(props.exprires),
			Issuer:    "exampleIssuer",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(props.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
