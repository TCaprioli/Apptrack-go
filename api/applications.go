package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	db "www.github.com/TCaprioli/Apptrack-go/db/sqlc"
)

func applicationHandler(store *db.Store, ctx context.Context) http.Handler {
	router := http.NewServeMux()
	router.Handle("POST /applications", handleFuncWithCtx(createApplication, store, ctx))
	router.Handle("GET /applications", handleFuncWithCtx(readAllApplications, store, ctx))

	return router
}
func applicationIdHandler(store *db.Store, ctx context.Context) http.Handler {
	router := http.NewServeMux()
	router.Handle("GET /applications/{id}", handleFuncWithCtx(readApplication, store, ctx))
	router.Handle("PUT /applications/{id}", handleFuncWithCtx(updateApplication, store, ctx))
	router.Handle("DELETE /applications/{id}", handleFuncWithCtx(deleteApplication, store, ctx))

	return router
}

func createApplication(w http.ResponseWriter, r *http.Request, store *db.Store, ctx context.Context) {
	var params db.CreateApplicationParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}
	jobApplication, err := store.CreateApplication(ctx, params)
	if err != nil {
		http.Error(w, "Failed to create job application", http.StatusInternalServerError)
		log.Printf("Error creating job application: %v", err)
		return
	}
	log.Printf("Creating job application...")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jobApplication); err != nil {
		http.Error(w, "Failed to encode job application response", http.StatusInternalServerError)
		log.Printf("Error encoding job application to response: %v", err)
	}
}

func readAllApplications(w http.ResponseWriter, r *http.Request, store *db.Store, ctx context.Context) {
	user := r.Context().Value(UserKey).(UserContext)
	limitStr := r.URL.Query().Get("limit")
	lastIDStr := r.URL.Query().Get("lastId")
	// defaults
	limit := int32(10)
	lastID := int32(0)

	if limitStr != "" {
		parsedLimit, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil {
			log.Println("Invalid limit, using default:", err)
		}
		limit = int32(parsedLimit)
	}
	if lastIDStr != "" {
		parsedId, err := strconv.ParseInt(lastIDStr, 10, 32)
		if err != nil {
			log.Println("Invalid id, using default:", err)
		}
		lastID = int32(parsedId)
	}

	params := db.ListApplicationsParams{
		UserID: user.id,
		ID:     lastID,
		Limit:  limit,
	}
	jobApplications, err := store.ListApplications(ctx, params)
	if err != nil {
		http.Error(w, "Failed to list job applications", http.StatusInternalServerError)
		return
	}

	// If there are no applications set it to an empty array
	// This is to prevent the empty list from being returned as `null`
	if len(jobApplications) == 0 {
		jobApplications = []db.Application{}
	}

	log.Printf("Fetching job applications...")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jobApplications); err != nil {
		http.Error(w, "Failed to encode job applications response", http.StatusInternalServerError)
	}

}

func readApplication(w http.ResponseWriter, r *http.Request, store *db.Store, ctx context.Context) {
	user := r.Context().Value(UserKey).(UserContext)

	idStr := r.PathValue("id")
	id := parseId(idStr, w)

	jobApplication, err := store.GetApplication(ctx, id)
	if err != nil {
		http.Error(w, "Failed to get job application data", http.StatusInternalServerError)
	}

	if jobApplication.UserID != user.id {
		http.Error(w, "You are not authorized to get this job application", http.StatusUnauthorized)

	}

	log.Printf("Fetching job id %v...", id)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jobApplication); err != nil {
		http.Error(w, "Failed to encode job applications response", http.StatusInternalServerError)
	}
}

func updateApplication(w http.ResponseWriter, r *http.Request, store *db.Store, ctx context.Context) {
	idStr := r.PathValue("id")
	id := parseId(idStr, w)
	var params db.UpdateApplicationParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Failed decoding request body", http.StatusBadRequest)
	}
	paramsWithId := db.UpdateApplicationParams{
		ID:              id,
		JobTitle:        params.JobTitle,
		Company:         params.Company,
		ApplicationDate: params.ApplicationDate,
		Status:          params.Status,
		Location:        params.Location,
		Notes:           params.Notes,
	}
	jobApplication, err := store.UpdateApplication(ctx, paramsWithId)
	if err != nil {
		http.Error(w, "Failed to update job application", http.StatusInternalServerError)
	}

	log.Printf("Updating job id %v...", id)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jobApplication); err != nil {
		http.Error(w, "Failed to encode job application response", http.StatusInternalServerError)
	}
}

func deleteApplication(w http.ResponseWriter, r *http.Request, store *db.Store, ctx context.Context) {
	idStr := r.PathValue("id")
	id := parseId(idStr, w)

	err := store.DeleteApplication(ctx, id)
	if err != nil {
		http.Error(w, "Failed to delete job application", http.StatusInternalServerError)
	}
	log.Printf("Deleting job id %v...", id)

}
