package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// helper

func fazerPedido(t *testing.T, metodo, caminho string, body io.Reader, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	app := setupRouter(time.Now())

	req := httptest.NewRequest(metodo, caminho, body)
	for chave, valor := range headers {
		req.Header.Set(chave, valor)
	}

	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	return rec
}

// handlers tests

func TestRootHandler(t *testing.T) {
	rec := fazerPedido(t, http.MethodGet, "/", nil, nil)

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
	app := setupRouter(time.Now().Add(-2 * time.Hour))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

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
	rec := fazerPedido(t, http.MethodGet, "/hello/Pedro", nil, nil)

	if rec.Code != http.StatusOK {
		t.Errorf("esperava status 200, recebi %d", rec.Code)
	}
	if rec.Body.String() != "Olá, Pedro!" {
		t.Errorf("esperava body %q, recebi %q", "Olá, Pedro!", rec.Body.String())
	}
}

func TestUserHandler(t *testing.T) {
	rec := fazerPedido(t, http.MethodGet, "/user/42", nil, nil)

	if rec.Code != http.StatusOK {
		t.Errorf("esperava status 200, recebi %d", rec.Code)
	}
	if rec.Body.String() != "User Profile: 42" {
		t.Errorf("esperava body %q, recebi %q", "User Profile: 42", rec.Body.String())
	}
}

func TestSearchHandlerComPagina(t *testing.T) {
	rec := fazerPedido(t, http.MethodGet, "/search?q=golang&page=2", nil, nil)

	if rec.Code != http.StatusOK {
		t.Errorf("esperava status 200, recebi %d", rec.Code)
	}
	if rec.Body.String() != "Buscando por 'golang' na página 2" {
		t.Errorf("body errado: %q", rec.Body.String())
	}
}

func TestSearchHandlerDefaultPage(t *testing.T) {
	rec := fazerPedido(t, http.MethodGet, "/search?q=golang", nil, nil)

	if rec.Code != http.StatusOK {
		t.Errorf("esperava status 200, recebi %d", rec.Code)
	}
	if rec.Body.String() != "Buscando por 'golang' na página 1" {
		t.Errorf("body errado: %q", rec.Body.String())
	}
}

func TestSearchHandlerSemQ(t *testing.T) {
	rec := fazerPedido(t, http.MethodGet, "/search", nil, nil)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("esperava status 400, recebi %d", rec.Code)
	}
}

func TestEchoHandlerSucesso(t *testing.T) {
	rec := fazerPedido(t, http.MethodPost, "/echo",
		strings.NewReader(`{"payload": "teste"}`),
		map[string]string{"X-App-Token": "secret123"})

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
	rec := fazerPedido(t, http.MethodPost, "/echo",
		strings.NewReader(""),
		map[string]string{"X-App-Token": "secret123"})

	if rec.Code != http.StatusBadRequest {
		t.Errorf("esperava status 400, recebi %d", rec.Code)
	}
}

func TestEchoHandlerJSONInvalido(t *testing.T) {
	rec := fazerPedido(t, http.MethodPost, "/echo",
		strings.NewReader("isto não é JSON"),
		map[string]string{"X-App-Token": "secret123"})

	if rec.Code != http.StatusBadRequest {
		t.Errorf("esperava status 400, recebi %d", rec.Code)
	}
}

// Middleware tests

func TestAuthMiddlewareSemToken(t *testing.T) {
	rec := fazerPedido(t, http.MethodPost, "/echo",
		strings.NewReader(`{"payload":"x"}`),
		nil)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("esperava 401, recebi %d", rec.Code)
	}
}

func TestAuthMiddlewareTokenErrado(t *testing.T) {
	rec := fazerPedido(t, http.MethodPost, "/echo",
		strings.NewReader(`{"payload":"x"}`),
		map[string]string{"X-App-Token": "errado"})

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("esperava 401, recebi %d", rec.Code)
	}
}

func TestAuthMiddlewareTokenCorreto(t *testing.T) {
	rec := fazerPedido(t, http.MethodPost, "/echo",
		strings.NewReader(`{"payload":"x"}`),
		map[string]string{"X-App-Token": "secret123"})

	if rec.Code != http.StatusOK {
		t.Errorf("esperava 200, recebi %d", rec.Code)
	}
}