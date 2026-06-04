package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"desafio3/services"
)

type Handlers struct {
	Service *services.AuthService
}

func NewHandlers(s *services.AuthService) *Handlers {
	return &Handlers{Service: s}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	user, err := h.Service.Register(req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserExists):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, services.ErrUsernameShort):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "erro interno", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}
	
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	token, err := h.Service.Login(req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorized):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, "erro interno", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}

func (h *Handlers) Profile(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("X-Session-Token")

	if token == "" {
		http.Error(w, "token em falta", http.StatusUnauthorized)
		return
	}

	user, err := h.Service.GetUserByToken(token)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorized):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, "erro interno", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("X-Session-Token")
	if token == "" {
		http.Error(w, "token em falta", http.StatusUnauthorized)
		return
	}

	err := h.Service.Logout(token)
	if err != nil {
		if errors.Is(err, services.ErrUnauthorized) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, "erro interno", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "sessão terminada com sucesso",
	})
}

func (h *Handlers) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("X-Session-Token")
	if token == "" {
		http.Error(w, "token em falta", http.StatusUnauthorized)
		return
	}

	users, err := h.Service.ListUsers(token)
	if err != nil {
		if errors.Is(err, services.ErrUnauthorized) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, "erro interno", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}