package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	db "www.github.com/TCaprioli/Apptrack-go/db/sqlc"
)

type UserResponse struct {
	ID    int32  `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func userHandler(store *db.Store, ctx context.Context) http.Handler {
	router := http.NewServeMux()
	router.Handle("POST /users/register", handleFuncWithCtx(registerUser, store, ctx))
	router.Handle("POST /users/login", handleFuncWithCtx(loginUser, store, ctx))
	router.Handle("POST /users/me", handleFuncWithCtx(verifyUser, store, ctx))
	router.Handle("DELETE /users/logout", http.HandlerFunc(logoutUser))
	return router
}

func registerUser(w http.ResponseWriter, r *http.Request, store *db.Store, ctx context.Context) {
	var params db.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	if len(params.Password) < MinimumPasswordLength || len(params.Password) > MaximumPasswordLength {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		log.Printf("Passwords must be between %v & %v", MinimumPasswordLength, MaximumPasswordLength)
		return
	}
	hashedPassword, hashErr := hashPassword(params.Password)
	if hashErr != nil {
		http.Error(w, hashErr.Error(), http.StatusBadRequest)
		log.Printf("Error creating user: %v", hashErr)
		return
	}

	paramsWithHash := db.CreateUserParams{
		Email:    params.Email,
		Name:     params.Name,
		Password: hashedPassword,
	}
	user, err := store.CreateUser(ctx, paramsWithHash)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		log.Printf("Error creating user: %v", err)
		return
	}
	tokenString, err := createToken(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Failed to sign token", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	userWithToken := UserResponse{Token: tokenString, ID: user.ID, Email: user.Email}

	log.Printf("Creating user...")

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userWithToken); err != nil {
		http.Error(w, "Failed to encode user response", http.StatusInternalServerError)
		log.Printf("Error encoding user to response: %v", err)
		return
	}
}

func loginUser(w http.ResponseWriter, r *http.Request, store *db.Store, ctx context.Context) {
	var params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	user, err := store.GetUser(ctx, params.Email)
	if err != nil {
		http.Error(w, "Invalid email", http.StatusUnauthorized)
		return
	}

	if err := checkPasswordHash(params.Password, user.Password); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	tokenString, err := createToken(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Failed to sign token", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	userWithToken := UserResponse{Token: tokenString, ID: user.ID, Email: user.Email}

	log.Printf("Logging in...")

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userWithToken); err != nil {
		http.Error(w, "Failed to encode users response", http.StatusInternalServerError)
		return
	}
}

func verifyUser(w http.ResponseWriter, r *http.Request, store *db.Store, ctx context.Context) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Unauthorized: missing cookie", http.StatusUnauthorized)
		return
	}

	token, err := verifyToken(cookie.Value)
	if err != nil {
		http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
		log.Printf("Token verification failed: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(token.Claims); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
func logoutUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("Logging out...")
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(-24 * time.Hour),
	})
	w.WriteHeader(http.StatusOK)
}
