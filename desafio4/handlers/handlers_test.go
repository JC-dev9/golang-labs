package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"desafio4/repository"
	"desafio4/services"
)

func setupTestHandlers(t *testing.T) *Handlers {
	t.Helper()
	
	// Verifica primeiro se já consegue ver a pasta "templates".
	// Só recua uma pasta se a "templates" não estiver visível.
	if _, err := os.Stat("templates"); os.IsNotExist(err) {
		_ = os.Chdir("..") 
	}

	userRepo := repository.NewMemoryUserRepository()
	sessionRepo := repository.NewMemorySessionRepository()
	service := services.NewAuthService(userRepo, sessionRepo)

	return NewHandlers(service)
}

func TestRegister(t *testing.T) {
	t.Run("GET - Mostra formulário de registo", func(t *testing.T) {
		h := setupTestHandlers(t)
		req := httptest.NewRequest(http.MethodGet, "/user/register", nil)
		rec := httptest.NewRecorder()

		h.RegisterGET(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("esperava 200 OK, recebi %d", rec.Code)
		}
		if !strings.Contains(rec.Body.String(), "<form action=\"/user/register\" method=\"POST\">") {
			t.Error("a página não contém o formulário correto")
		}
	})

	t.Run("POST - Sucesso redireciona para login", func(t *testing.T) {
		h := setupTestHandlers(t)
		body := strings.NewReader("username=intern1&password=safe-password")
		req := httptest.NewRequest(http.MethodPost, "/user/register", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.RegisterPOST(rec, req)

		if rec.Code != http.StatusSeeOther {
			t.Fatalf("esperava redirecionamento (303), recebi %d", rec.Code)
		}
		if rec.Header().Get("Location") != "/user/login" {
			t.Errorf("não redirecionou para o login, foi para: %s", rec.Header().Get("Location"))
		}
	})

	t.Run("POST - Username já existe mostra erro na página", func(t *testing.T) {
		h := setupTestHandlers(t)
		_, _ = h.Service.Register("intern1", "safe-password")

		body := strings.NewReader("username=intern1&password=outra-password")
		req := httptest.NewRequest(http.MethodPost, "/user/register", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.RegisterPOST(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("esperava 200 OK, recebi %d", rec.Code)
		}
		if !strings.Contains(rec.Body.String(), "username já existe") {
			t.Error("a página não mostrou a mensagem de erro correta")
		}
	})
}

func TestLogin(t *testing.T) {
	t.Run("GET - Mostra formulário de login", func(t *testing.T) {
		h := setupTestHandlers(t)
		req := httptest.NewRequest(http.MethodGet, "/user/login", nil)
		rec := httptest.NewRecorder()

		h.LoginGET(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("esperava 200 OK, recebi %d", rec.Code)
		}
	})

	t.Run("POST - Sucesso redireciona para perfil", func(t *testing.T) {
		h := setupTestHandlers(t)
		_, _ = h.Service.Register("intern1", "safe-password")

		body := strings.NewReader("username=intern1&password=safe-password")
		req := httptest.NewRequest(http.MethodPost, "/user/login", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.LoginPOST(rec, req)

		if rec.Code != http.StatusSeeOther {
			t.Fatalf("esperava redirecionamento (303), recebi %d", rec.Code)
		}
		if !strings.HasPrefix(rec.Header().Get("Location"), "/user/profile?token=") {
			t.Errorf("não redirecionou para o perfil com o token. Location: %s", rec.Header().Get("Location"))
		}
	})

	t.Run("POST - Credenciais inválidas mostram erro", func(t *testing.T) {
		h := setupTestHandlers(t)
		_, _ = h.Service.Register("intern1", "safe-password")

		body := strings.NewReader("username=intern1&password=errada")
		req := httptest.NewRequest(http.MethodPost, "/user/login", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.LoginPOST(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("esperava 200 OK (re-renderização), recebi %d", rec.Code)
		}
		if !strings.Contains(rec.Body.String(), "Credenciais inválidas") {
			t.Error("a página não mostrou a mensagem de erro")
		}
	})
}

func TestProfile(t *testing.T) {
	t.Run("GET - Sucesso renderiza perfil", func(t *testing.T) {
		h := setupTestHandlers(t)
		_, _ = h.Service.Register("intern1", "safe-password")
		token, _ := h.Service.Login("intern1", "safe-password")

		req := httptest.NewRequest(http.MethodGet, "/user/profile?token="+token, nil)
		rec := httptest.NewRecorder()

		h.Profile(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("esperava 200 OK, recebi %d", rec.Code)
		}
		if !strings.Contains(rec.Body.String(), "Bem-vindo, intern1!") {
			t.Error("a página não renderizou o nome do utilizador")
		}
	})

	t.Run("GET - Sem token redireciona para login", func(t *testing.T) {
		h := setupTestHandlers(t)
		req := httptest.NewRequest(http.MethodGet, "/user/profile", nil)
		rec := httptest.NewRecorder()

		h.Profile(rec, req)

		if rec.Code != http.StatusSeeOther {
			t.Fatalf("esperava redirecionamento (303), recebi %d", rec.Code)
		}
		if rec.Header().Get("Location") != "/user/login" {
			t.Error("não redirecionou para a página de login")
		}
	})
}

func TestLogout(t *testing.T) {
	t.Run("POST - Apaga sessão e redireciona", func(t *testing.T) {
		h := setupTestHandlers(t)
		_, _ = h.Service.Register("intern1", "safe-password")
		token, _ := h.Service.Login("intern1", "safe-password")

		// Simula o envio do token no corpo da requisição, como se fosse um formulário de logout
		body := strings.NewReader("token=" + token)
		req := httptest.NewRequest(http.MethodPost, "/user/logout", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.Logout(rec, req)

		if rec.Code != http.StatusSeeOther {
			t.Fatalf("esperava redirecionamento (303), recebi %d", rec.Code)
		}
		
		// Verifica se a sessão foi realmente apagada
		_, err := h.Service.GetUserByToken(token)
		if err == nil {
			t.Error("o logout não apagou a sessão do repositório")
		}
	})
}