package api

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/o1egl/paseto"
	"www.github.com/TCaprioli/Apptrack-go/utils"
)

var authorizationTypeBearer = "bearer"

type contextKey string

const UserKey contextKey = "user"

type UserContext struct {
	id int32
}

func authMiddleware(next http.Handler) http.Handler {
	utils.LoadEnv()
	symmetricKey := os.Getenv("SECRET_KEY")
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
		err := paseto.NewV2().Decrypt(token, []byte(symmetricKey), &user.id, nil)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
