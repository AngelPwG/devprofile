package router

import (
	"net/http"

	"github.com/AngelPwG/devprofile/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewRouter(h *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}))

	r.Post("/profiles", h.CreateProfile)
	r.Get("/profiles", h.GetProfiles)
	r.Get("/profiles/{username}", h.GetProfile)
	r.Put("/profiles/{username}", h.UpdateProfile)
	r.Delete("/profiles/{username}", h.DeleteProfile)
	r.Get("/audit", h.GetAuditLogs)

	return r
}
