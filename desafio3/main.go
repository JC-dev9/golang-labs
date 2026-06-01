package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"desafio3/handlers"
	"desafio3/repository"
	"desafio3/services"
)

func setupRouter(h *handlers.Handlers) http.Handler {
	r := chi.NewRouter()

	r.Post("/api/user/register", h.Register)
	r.Post("/api/user/login", h.Login)
	r.Get("/api/user/profile", h.Profile)

	return r
}

func main() {
	
	// No teu main.go, troca o services.NewAuthService() por:
	userRepo := repository.NewMemoryUserRepository()
	sessionRepo := repository.NewMemorySessionRepository()
	authService := services.NewAuthService(userRepo, sessionRepo)
	
	// bootstrap do utilizador admin
	if _, err := authService.Register("admin", "password123"); err != nil {
		panic(err)
	}

	h := handlers.NewHandlers(authService)
	app := setupRouter(h)
	http.ListenAndServe(":8080", app)
}