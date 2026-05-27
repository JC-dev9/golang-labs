package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)
type HealthResponse struct {
    Status string `json:"status"`
    Uptime string `json:"uptime"`
}

// handlers

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Ola Mundo!"))
}

func makeHealthHandler(startTime time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := HealthResponse{
			Status: "up",
			Uptime: time.Since(startTime).String(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Olá, " + name + "!"))
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("User Profile: " + id))
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	
	if q == "" {
		http.Error(w, "parametro 'q' em falta", http.StatusBadRequest)
		return
	}

	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			page = n
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Buscando por '%s' na página %d", q, page)
}


type EchoRequest struct {
	Payload string `json:"payload"`
}

type EchoResponse struct {
	Payload     string `json:"payload"`
	ProcessedAt string `json:"processed_at"`
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "body vazio", http.StatusBadRequest)
		return
	}

	var req EchoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if req.Payload == "" {
		http.Error(w, "payload vazio", http.StatusBadRequest)
		return
	}

	resp := EchoResponse{
		Payload:     req.Payload,
		ProcessedAt: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// middlewares

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		latency := time.Since(start)
		fmt.Printf("[%s] - %s - %s\n", r.Method, r.URL.Path, latency)
	})
}


func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-App-Token")
		if token != "secret123" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Unauthorized: token inválido ou ausente"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// main

func main() {
	startTime := time.Now()

	app := chi.NewRouter()

	app.Use(loggingMiddleware)

	app.Get("/", rootHandler)
	app.Get("/health", makeHealthHandler(startTime))
	app.Get("/hello/{name}", helloHandler)
	app.Get("/user/{id:[0-9]+}", userHandler)
	app.Get("/search", searchHandler)

	app.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/echo", echoHandler)
	})

	http.ListenAndServe(":8080", app)
}