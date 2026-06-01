package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"desafio3/repository"
	"desafio3/services"
)

func setupTestHandlers(t *testing.T) *Handlers {
	t.Helper()
	userRepo := repository.NewMemoryUserRepository()
	sessionRepo := repository.NewMemorySessionRepository()
	service := services.NewAuthService(userRepo, sessionRepo)

	return NewHandlers(service)
}

func TestRegisterSucesso(t *testing.T) {
	h := setupTestHandlers(t)

	body := strings.NewReader(`{"username":"intern1","password":"safe-password"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", body)
	rec := httptest.NewRecorder()
	h.Register(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("esperava 201, recebi %d", rec.Code)
	}

	var resp struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		CreatedAt string `json:"created_at"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("resposta não é JSON válido: %v", err)
	}

	if resp.Username != "intern1" {
		t.Errorf("esperava username %q, recebi %q", "intern1", resp.Username)

	}
	if len(resp.ID) != 8 {
		t.Errorf("esperava ID com 8 caracteres, recebi %d", len(resp.ID))
	}
	if resp.CreatedAt == "" {
		t.Error("created_at está vazio")
	}
}

func TestRegisterDuplicado(t *testing.T) {
	h := setupTestHandlers(t)

	if _, err := h.Service.Register("intern1", "safe-password"); err != nil {
		t.Fatalf("seed falhou: %v", err)
	}

	body := strings.NewReader(`{"username":"intern1","password":"outra-password"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", body)
	rec := httptest.NewRecorder()
	h.Register(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("esperava 409, recebi %d", rec.Code)
	}
}

func TestRegisterUsernameCurto(t *testing.T) {
	h := setupTestHandlers(t)

	body := strings.NewReader(`{"username":"ab","password":"safe-password"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", body)
	rec := httptest.NewRecorder()
	h.Register(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("esperava 400, recebi %d", rec.Code)
	}
}

func TestRegisterJSONInvalido(t *testing.T) {
	h := setupTestHandlers(t)

	body := strings.NewReader(`{"username": "intern1", "password":`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", body)
	rec := httptest.NewRecorder()
	h.Register(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("esperava 400, recebi %d", rec.Code)
	}
}

func TestLoginSucesso(t *testing.T) {
	h := setupTestHandlers(t)

	if _, err := h.Service.Register("intern1", "safe-password"); err != nil {
		t.Fatalf("seed falhou: %v", err)
	}

	body := strings.NewReader(`{"username":"intern1","password":"safe-password"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", body)
	rec := httptest.NewRecorder()
	h.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperava 200, recebi %d", rec.Code)
	}

	var resp struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("resposta não é JSON válido: %v", err)
	}

	if resp.Token == "" {
		t.Error("token está vazio")
	}
}

func TestLoginCredenciaisInvalidas(t *testing.T) {
	h := setupTestHandlers(t)

	if _, err := h.Service.Register("intern1", "safe-password"); err != nil {
		t.Fatalf("seed falhou: %v", err)
	}

	body := strings.NewReader(`{"username":"intern1","password":"errada"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", body)
	rec := httptest.NewRecorder()
	h.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("esperava 401, recebi %d", rec.Code)
	}
}

func TestProfileSucesso(t *testing.T) {
	h := setupTestHandlers(t)

	if _, err := h.Service.Register("intern1", "safe-password"); err != nil {
		t.Fatalf("seed falhou: %v", err)
	}
	token, err := h.Service.Login("intern1", "safe-password")
	if err != nil {
		t.Fatalf("login seed falhou: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/user/profile", nil)
	req.Header.Set("X-Session-Token", token)
	rec := httptest.NewRecorder()
	h.Profile(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperava 200, recebi %d", rec.Code)
	}

	var resp struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		CreatedAt string `json:"created_at"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("resposta não é JSON válido: %v", err)
	}

	if resp.Username != "intern1" {
		t.Errorf("esperava username %q, recebi %q", "intern1", resp.Username)
	}
}

func TestProfileSemToken(t *testing.T) {
	h := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodGet, "/api/user/profile", nil)
	rec := httptest.NewRecorder()
	h.Profile(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("esperava 401, recebi %d", rec.Code)
	}
}

func TestProfileTokenInvalido(t *testing.T) {
	h := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodGet, "/api/user/profile", nil)
	req.Header.Set("X-Session-Token", "token-falso")
	rec := httptest.NewRecorder()
	h.Profile(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("esperava 401, recebi %d", rec.Code)
	}
}