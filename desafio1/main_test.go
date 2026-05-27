package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"encoding/json"
	"strings"

	"github.com/go-chi/chi/v5"
)


func TestRootHandler(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	rootHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("esperava status 200, recebi %d", rec.Code)
	}

	if ct := rec.Header().Get("Content-Type"); ct != "text/plain" {
    t.Errorf("esperava Content-Type %q, recebi %q", "text/plain", ct)	
	}


	if rec.Body.String() != "Ola Mundo!" {
		t.Errorf("esperava body %q, recebi %q", "Ola Mundo!", rec.Body.String())
	}

}

func TestHealthHandler(t *testing.T) {
	startTime := time.Now().Add(-2 * time.Hour)
	handler := makeHealthHandler(startTime)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("esperava status 200, recebi %d", rec.Code)
	}

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("esperava Content-Type %q, recebi %q", "application/json", ct)
	}

	var body struct {
		Status string `json:"status"`
		Uptime string `json:"uptime"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("resposta não é JSON válido: %v", err)
	}

	if body.Status != "up" {
		t.Errorf("esperava status %q, recebi %q", "up", body.Status)
	}

	uptime, err := time.ParseDuration(body.Uptime)
	if err != nil {
		t.Fatalf("uptime %q não é uma duração válida: %v", body.Uptime, err)
	}

	if uptime < 2*time.Hour || uptime > 2*time.Hour+time.Second {
		t.Errorf("esperava uptime ~2h, recebi %v", uptime)
	}
}

func TestHello(t *testing.T) {


	req := httptest.NewRequest(http.MethodGet, "/hello/Pedro", nil)
	rec := httptest.NewRecorder()

	app := chi.NewRouter()
	app.Get("/hello/{name}", helloHandler)
	app.ServeHTTP(rec,req)

	if rec.Code != http.StatusOK {
		t.Errorf("esperava status 200, recebi %d", rec.Code)
	}

	if rec.Body.String() != "Olá, Pedro!" {
		t.Errorf("esperava body %q, recebi %q", "Olá, Pedro!", rec.Body.String())
	}
}


func TestUserHandler(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/user/42", nil)
	rec := httptest.NewRecorder()

	app := chi.NewRouter()
	app.Get("/user/{id:[0-9]+}", userHandler)
	app.ServeHTTP(rec,req)

	if rec.Code != http.StatusOK {
		t.Errorf("esperava status 200, recebi %d", rec.Code)
	}

	if rec.Body.String() != "User Profile: 42" {
		t.Errorf("esperava body %q, recebi %q", "User Profile: 42", rec.Body.String())
	}
}

func TestSearchHandlerComPagina(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/search?q=golang&page=2", nil)
	rec := httptest.NewRecorder()

	searchHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("esperava status 200, recebi %d", rec.Code)
	}
	if rec.Body.String() != "Buscando por 'golang' na página 2" {
		t.Errorf("body errado: %q", rec.Body.String())
	}
}

func TestSearchHandlerDefaultPage(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/search?q=golang", nil)
	rec := httptest.NewRecorder()

	searchHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("esperava status 200, recebi %d", rec.Code)
	}
	if rec.Body.String() != "Buscando por 'golang' na página 1" {
		t.Errorf("body errado: %q", rec.Body.String())
	}
}

func TestSearchHandlerSemQ(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/search", nil)
	rec := httptest.NewRecorder()

	searchHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("esperava status 400, recebi %d", rec.Code)
	}
}

func TestEchoHandlerSucesso(t *testing.T) {
	body := strings.NewReader(`{"payload": "teste"}`)
	req := httptest.NewRequest(http.MethodPost, "/echo", body)
	rec := httptest.NewRecorder()

	echoHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("esperava status 200, recebi %d", rec.Code)
	}

	var resp struct {
		Payload     string `json:"payload"`
		ProcessedAt string `json:"processed_at"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("resposta não é JSON válido: %v", err)
	}

	if resp.Payload != "teste" {
		t.Errorf("payload errado: %q", resp.Payload)
	}
	if resp.ProcessedAt == "" {
		t.Errorf("processed_at vazio")
	}
}

func TestEchoHandlerBodyVazio(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(""))
	rec := httptest.NewRecorder()

	echoHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("esperava status 400, recebi %d", rec.Code)
	}
}

func TestEchoHandlerJSONInvalido(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader("isto não é JSON"))
	rec := httptest.NewRecorder()

	echoHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("esperava status 400, recebi %d", rec.Code)
	}
}

// Testes dos middlewares de autenticação

func TestAuthMiddlewareSemToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(`{"payload":"x"}`))
	rec := httptest.NewRecorder()

	app := chi.NewRouter()
	app.Use(authMiddleware)
	app.Post("/echo", echoHandler)
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("esperava 401, recebi %d", rec.Code)
	}
}

func TestAuthMiddlewareTokenErrado(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(`{"payload":"x"}`))
	req.Header.Set("X-App-Token", "errado")
	rec := httptest.NewRecorder()

	app := chi.NewRouter()
	app.Use(authMiddleware)
	app.Post("/echo", echoHandler)
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("esperava 401, recebi %d", rec.Code)
	}
}

func TestAuthMiddlewareTokenCorreto(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(`{"payload":"x"}`))
	req.Header.Set("X-App-Token", "secret123")
	rec := httptest.NewRecorder()

	app := chi.NewRouter()
	app.Use(authMiddleware)
	app.Post("/echo", echoHandler)
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("esperava 200, recebi %d", rec.Code)
	}
}