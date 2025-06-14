package controller

import (
	"context"
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
)


type contextKey string
const userContextKey = contextKey("firebase_uid")

func AuthMiddleware(authClient *auth.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid Authorization header format. Expected 'Bearer <token>'", http.StatusUnauthorized)
				return
			}
			idToken := parts[1]

			token, err := authClient.VerifyIDToken(r.Context(), idToken)
			if err != nil {
				log.Printf("error verifying ID token: %v\n", err)
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, token.UID)


			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}