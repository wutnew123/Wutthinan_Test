package middleware

import (
	"back-end/utils"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey   contextKey = "userID"
	UsernameKey contextKey = "username"
)

var JwtKey = []byte("your_secret_key")

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if !strings.HasPrefix(auth, "Bearer ") {
			utils.RespondJSON(w, http.StatusUnauthorized, "error", "Missing or invalid Authorization header", nil)
			return
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return JwtKey, nil
		})

		if err != nil || !token.Valid {
			utils.RespondJSON(w, http.StatusUnauthorized, "error", "Invalid token", nil)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.RespondJSON(w, http.StatusUnauthorized, "error", "Invalid token claims", nil)
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			utils.RespondJSON(w, http.StatusUnauthorized, "error", "Invalid user_id", nil)
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			utils.RespondJSON(w, http.StatusUnauthorized, "error", "Invalid username", nil)
			return
		}

		userID := int(userIDFloat)

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UsernameKey, username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
