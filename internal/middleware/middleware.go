package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Zrossiz/gophermart/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDContextKey   contextKey = "userID"
	UserNameContextKey contextKey = "userName"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("accesstoken")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value

		secretKey := []byte(config.AppConfig.AccessTokenSecret)

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			var userID int
			var okID bool
			if idVal, ok := claims["userId"]; ok {
				switch v := idVal.(type) {
				case string:
					userID, err = strconv.Atoi(v)
					if err != nil {
						http.Error(w, "unauthorized", http.StatusUnauthorized)
						return
					}
					okID = true
				}
			}

			userName, okName := claims["userName"].(string)

			if !okID || !okName {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
			ctx = context.WithValue(ctx, UserNameContextKey, userName)

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
