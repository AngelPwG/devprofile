package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	cache "github.com/AngelPwG/devprofile/internal/cache"
	db "github.com/AngelPwG/devprofile/internal/db"
	builder "github.com/AngelPwG/devprofile/internal/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	db *db.DB
}

func NewHandler(db *db.DB) *Handler {
	return &Handler{db: db}
}

func jsonError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (h *Handler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var username string
	if err := json.NewDecoder(r.Body).Decode(&username); err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	profile, repos, err := builder.BuildProfile(username)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.db.InsertProfile(*profile); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.db.InsertRepositories(repos, profile.ID); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	ip := r.RemoteAddr
	if err := h.db.InsertAuditLog("CREATE", username, ip); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(profile)
}

func (h *Handler) GetProfiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	profiles, err := h.db.GetProfiles()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	json.NewEncoder(w).Encode(profiles)
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	username := chi.URLParam(r, "username")
	if username == "" {
		jsonError(w, http.StatusBadRequest, "username is required")
		return
	}
	profile, err := h.db.GetProfile(username)
	if err == sql.ErrNoRows {
		jsonError(w, http.StatusNotFound, "profile not found")
		return
	}
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	json.NewEncoder(w).Encode(profile)
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	username := chi.URLParam(r, "username")
	if username == "" {
		jsonError(w, http.StatusBadRequest, "username is required")
		return
	}
	oldProfile, err := h.db.GetProfile(username)
	if err == sql.ErrNoRows {
		jsonError(w, http.StatusNotFound, "profile not found")
		return
	}
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	canRefresh, remaining, err := cache.CanRefresh(oldProfile.UpdatedAt)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !canRefresh {
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]any{
			"error":               "too_soon",
			"retry_after_seconds": remaining,
		})
		return
	}

	profile, repos, err := builder.BuildProfile(username)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.db.UpdateProfile(*profile); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.db.DeleteRepositories(profile.ID); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.db.InsertRepositories(repos, profile.ID); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	ip := r.RemoteAddr
	if err := h.db.InsertAuditLog("UPDATE", username, ip); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	json.NewEncoder(w).Encode(profile)
}

func (h *Handler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	username := chi.URLParam(r, "username")
	if username == "" {
		jsonError(w, http.StatusBadRequest, "username is required")
		return
	}
	_, err := h.db.GetProfile(username)
	if err == sql.ErrNoRows {
		jsonError(w, http.StatusNotFound, "profile not found")
		return
	}
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.db.DeleteProfile(username); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	ip := r.RemoteAddr
	if err := h.db.InsertAuditLog("DELETE", username, ip); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "profile deleted successfully"})
}
