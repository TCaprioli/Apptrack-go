package api

import (
	"context"
	"net/http"
	"strings"
)

var authorizationTypeBearer = "bearer"

type contextKey string

const UserKey contextKey = "user"

type UserContext struct {
	ID int32
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fields := strings.Fields(header)
		if len(fields) < 2 {
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			http.Error(w, "unsupported authorization type", http.StatusUnauthorized)
			return
		}

		token := fields[1]
		var user UserContext
		verified, err := verifyToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		claims, ok:= verified.Claims.(*MyClaims)
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
        if origin == "http://localhost:5173" || origin == "https://apptrack-bice.vercel.app" {
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
