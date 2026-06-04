package services

import (
	"desafio4/repository"
	"errors"
	"testing"
	"time"
)

// helper para os testes (função auxiliar)
func setupTestService(t *testing.T) *AuthService {
	t.Helper() 

	userRepo := repository.NewMemoryUserRepository()
	sessionRepo := repository.NewMemorySessionRepository()

	return NewAuthService(userRepo, sessionRepo)
}

func TestRegisterUsernameShort(t *testing.T) {

	service := setupTestService(t)
	
	_,err := service.Register("abc", "password123")

	if err == nil {
		t.Fatal("esperava erro para username curto, mas não recebi")
	}

	if !errors.Is(err, ErrUsernameShort) {
		t.Errorf("esperava ErrUsernameShort, mas recebi: %v", err)
	}
}

func TestRegisterUsernameValido(t *testing.T) {
	service := setupTestService(t)
	
	_, err := service.Register("juan", "password123")
	if err != nil {
        t.Fatalf("primeiro Register não devia falhar: %v", err)
    }

	_, err = service.Register("juan", "outrapassword")
	if err == nil {
		t.Fatal("esperava erro para username já existente, mas não recebi")
	}

	if !errors.Is(err, ErrUserExists) {
		t.Errorf("esperava ErrUserExists, recebi: %v", err)
	}
}

func TestRegisterSucesso(t *testing.T) {
	service := setupTestService(t)
	antes := time.Now()

	user, err := service.Register("juan", "password123")
	if err != nil {
		t.Fatalf("não esperava erro, recebi: %v", err)
	}

	depois := time.Now()

	if user.Username != "juan" {
		t.Errorf("esperava username %q, recebi %q", "juan", user.Username)
	}

	if user.ID == "" {
		t.Error("ID está vazio")
	}
	if len(user.ID) != 8 {
		t.Errorf("esperava ID com 8 caracteres, recebi %d (%q)", len(user.ID), user.ID)
	}

	if user.CreatedAt.Before(antes) || user.CreatedAt.After(depois) {
		t.Errorf("CreatedAt %v fora do intervalo [%v, %v]", user.CreatedAt, antes, depois)
	}
}

// login tests

func TestLoginUsernameInexistente(t *testing.T) {

	service := setupTestService(t)

	_, err := service.Login("inexistente", "password123")

	if err == nil {
		t.Fatal("esperava erro para username inexistente, mas não recebi")
	}

	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("esperava ErrUnauthorized, mas recebi: %v", err)
	}
}

func TestLoginPasswordErrada(t *testing.T) {
	service := setupTestService(t)

	_, err := service.Register("juan", "password123")
	if err != nil {
		t.Fatalf("primeiro Register não devia falhar: %v", err)
	}

	_, err = service.Login("juan", "errada")
	if err == nil {
		t.Fatal("esperava erro para password errada, mas não recebi")
	}

	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("esperava ErrUnauthorized, mas recebi: %v", err)
	}
}

func TestLoginSucesso(t *testing.T) {
	service := setupTestService(t)
	_, err := service.Register("juan", "password123")

	if err != nil {
		t.Fatalf("primeiro Register não devia falhar: %v", err)
	}

	token, err := service.Login("juan", "password123")
	if err != nil {
		t.Fatalf("Login não devia falhar, mas recebi: %v", err)
	}
	
	if token == "" {
		t.Error("token está vazio")
	}
}

func TestTokenInvalido(t *testing.T) {
    service := setupTestService(t)

    _, err := service.GetUserByToken("token-invalido")
    if err == nil {
        t.Fatal("esperava erro para token inválido, mas não recebi")
    }

    if !errors.Is(err, ErrUnauthorized) {
        t.Errorf("esperava ErrUnauthorized, mas recebi: %v", err)
    }
}

func TestTokenValido(t *testing.T) {
	service := setupTestService(t)
	_, err := service.Register("juan", "password123")
	if err != nil {
		t.Fatalf("primeiro Register não devia falhar: %v", err)
	}

	token, err := service.Login("juan", "password123")
	if err != nil {
		t.Fatalf("Login não devia falhar, mas recebi: %v", err)
	}

	user, err := service.GetUserByToken(token)
	if err != nil {
		t.Fatalf("GetUserByToken não devia falhar, mas recebi: %v", err)
	}

	if user.Username != "juan" {
		t.Errorf("esperava username %q, recebi %q", "juan", user.Username)
	}
}

func TestLogoutSucesso(t *testing.T) {
	service := setupTestService(t)
	_, _ = service.Register("juan", "password123")
	token, _ := service.Login("juan", "password123")

	err := service.Logout(token)
	if err != nil {
		t.Fatalf("Logout falhou: %v", err)
	}

	_, err = service.GetUserByToken(token)
	if err == nil {
		t.Fatal("Esperava erro ao usar token apagado, não recebi")
	}
}

func TestLogoutTokenInvalido(t *testing.T) {
	service := setupTestService(t)
	err := service.Logout("token-inventado")
	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("Esperava ErrUnauthorized, recebi: %v", err)
	}
}

func TestListUsersSucesso(t *testing.T) {
	service := setupTestService(t)
	_, _ = service.Register("admin", "pass1")
	_, _ = service.Register("user2", "pass2")
	token, _ := service.Login("admin", "pass1")

	users, err := service.ListUsers(token)
	if err != nil {
		t.Fatalf("ListUsers falhou: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Esperava 2 utilizadores, recebi %d", len(users))
	}
}

func TestListUsersTokenInvalido(t *testing.T) {
	service := setupTestService(t)
	_, err := service.ListUsers("token-falso")
	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("Esperava ErrUnauthorized, recebi: %v", err)
	}
}