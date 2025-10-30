package api

import (
	"context"
	"log"
	"net/http"
)

var authorizationTypeBearer = "bearer"

type contextKey string

const UserKey contextKey = "user"

type UserContext struct {
	ID int32
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			log.Printf(": %v", err)
			http.Error(w, "Did not find valid cookie", http.StatusUnauthorized)
			return
		}
		tokenString := c.Value
		var user UserContext
		verified, err := verifyToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		claims, ok := verified.Claims.(*MyClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
		user = UserContext{ID: int32(claims.UserID)}
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:5173" || origin == "http://localhost:3000" || origin == "https://apptrack-bice.vercel.app" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
