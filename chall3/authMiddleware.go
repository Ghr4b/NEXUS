package main

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDKey contextKey = "user_id"
const userRoleKey contextKey = "user_role"

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the cookie from the request
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Parse and validate the JWT token
		tokenStr := cookie.Value
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract user_id and inject into context
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID, ok := claims["user_id"].(string)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			role, ok := claims["role"].(float64)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			ctx = context.WithValue(ctx, userRoleKey, int(role))

			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	})
}

func adminMiddleware(next http.Handler) http.Handler {
	return authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(userRoleKey).(int)
		if !ok || role != 0 { // change this to role = 0
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}))
}
