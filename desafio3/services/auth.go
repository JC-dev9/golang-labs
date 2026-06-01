package services

import (
	"crypto/rand"
	"desafio3/repository"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameShort = errors.New("username deve ter pelo menos 4 caracteres")
	ErrUserExists    = errors.New("username já existe")
	ErrUnauthorized   = errors.New("credenciais inválidas")
)

type AuthService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
}

func NewAuthService(ur repository.UserRepository, sr repository.SessionRepository) *AuthService {
	return &AuthService{
		userRepo:    ur,
		sessionRepo: sr,
	}
}

func generateRandomHex(bytes int) string {
	b := make([]byte, bytes)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func generateID() string {
	return generateRandomHex(4) // 8 caracteres hexadecimais
}

func (s *AuthService) Register(username, password string) (repository.User, error) {
	if len(username) < 4 {
		return repository.User{}, ErrUsernameShort
	}

	_, err := s.userRepo.FindByUsername(username)
	if err == nil {
		return repository.User{}, ErrUserExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return repository.User{}, err
	}

	user := repository.User{
		ID:        generateID(),
		Username:  username,
		Password:  string(hash),
		CreatedAt: time.Now(),
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return repository.User{}, err
	}

	return user, nil
}

func (s *AuthService) Login(username, password string) (string, error) {
	
	// Buscar o user pelo username
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", ErrUnauthorized
	}

	// A password em hash agora está DENTRO do user:
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrUnauthorized
	}

	token := generateRandomHex(16)
	
	// Criar e guardar a sessão com o repositório:
	novaSessao := repository.Session{
		Token:     token,
		UserID:    user.ID,
		CreatedAt: time.Now(),
	}
	
	err = s.sessionRepo.Create(novaSessao)
	if err != nil {
	    return "", err
	}

	return token, nil
}

func (s *AuthService) GetUserByToken(token string) (repository.User, error) {
	// Buscar a sessão pelo token:
	sessao, err := s.sessionRepo.FindByToken(token)
	if err != nil { // Se não encontrar a sessão, o token é inválido
		return repository.User{}, ErrUnauthorized
	}

	// Se a sessão é válida, vai usar o UserID da sessão para buscar o utilizador:
	user, err := s.userRepo.FindByID(sessao.UserID)
	if err != nil {
		return repository.User{}, ErrUnauthorized
	}

	return user, nil
}