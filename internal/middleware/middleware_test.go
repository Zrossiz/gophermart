package middleware

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/Zrossiz/gophermart/internal/config"
	"github.com/Zrossiz/gophermart/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestJWTMiddleware(t *testing.T) {
	config.AppConfig = &config.Config{AccessTokenSecret: "testsecret"}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(UserIDContextKey).(int)
		userName := r.Context().Value(UserNameContextKey).(string)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User ID: " + strconv.Itoa(userID) + ", User Name: " + userName))
	})

	middleware := JWTMiddleware(handler)

	t.Run("No token provided", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "unauthorized")
	})

	t.Run("Invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		req.AddCookie(&http.Cookie{Name: "accesstoken", Value: "invalidtoken"})

		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "unauthorized")
	})

	t.Run("Valid token", func(t *testing.T) {
		token, err := utils.GenerateJWT(utils.GenerateJWTProps{
			Secret:   []byte(config.AppConfig.AccessTokenSecret),
			Exprires: time.Now().Add(time.Hour),
			UserID:   123,
			Username: "testuser",
		})

		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		req.AddCookie(&http.Cookie{Name: "accesstoken", Value: token})

		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "User ID: 123, User Name: testuser")
	})
}
