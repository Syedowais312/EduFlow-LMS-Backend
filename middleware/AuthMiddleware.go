package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"goeduflow/utils"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey   contextKey = "userID"
	SchoolKey   contextKey = "schoolName"
	UsernameKey contextKey = "username"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return utils.JwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// 🔹 DEBUG: print parsed claims
		fmt.Printf("DEBUG parsed claims: %+v\n", claims)

		// extract correctly using the exact keys from GenerateJWT
		userIDf, ok1 := claims["user_id"].(float64)
		school, ok2 := claims["schoolName"].(string)
		username, ok3 := claims["username"].(string)

		if !ok1 || !ok2 || !ok3 {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, int(userIDf))
		ctx = context.WithValue(ctx, SchoolKey, school)
		ctx = context.WithValue(ctx, UsernameKey, username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
