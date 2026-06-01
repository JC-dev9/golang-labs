package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"desafio2/handlers"
	"desafio2/services"
)

func setupRouter(h *handlers.Handlers) http.Handler {
	r := chi.NewRouter()

	r.Post("/api/user/register", h.Register)
	r.Post("/api/user/login", h.Login)
	r.Get("/api/user/profile", h.Profile)

	return r
}

func main() {
	service := services.NewAuthService()

	// bootstrap do utilizador admin
	if _, err := service.Register("admin", "password123"); err != nil {
		panic(err)
	}

	h := handlers.NewHandlers(service)
	app := setupRouter(h)
	http.ListenAndServe(":8080", app)
}