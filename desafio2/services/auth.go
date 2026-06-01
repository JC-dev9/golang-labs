package services

import (
	"crypto/rand"
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
	users map[string]User
	passwords map[string][]byte // hash da password
	sessions map[string]string // token -> username
}

func NewAuthService() *AuthService {
	return &AuthService{
		users: make(map[string]User),
		passwords: make(map[string][]byte),
		sessions: make(map[string]string),
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

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *AuthService) Register(username, password string) (User, error) {
	if len(username) < 4 {
		return User{}, ErrUsernameShort
	}
	
	if _, existe := s.users[username]; existe {
		return User{}, ErrUserExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	
	if err != nil {
		return User{}, err
	}

	user := User{
		ID:        generateID(),
		Username:  username,
		CreatedAt: time.Now(),
	}
	s.passwords[username] = hash
	s.users[username] = user
	return user, nil
}

func (s *AuthService) Login(username, password string) (string, error) {
	_, existe := s.users[username]
	if !existe {
		return "", ErrUnauthorized
	}

	hash := s.passwords[username]
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		return "", ErrUnauthorized
	}

	token := generateRandomHex(16)
	s.sessions[token] = username
	return token, nil
}

func (s *AuthService) GetUserByToken(token string) (User, error) {
	username, existe := s.sessions[token]
	if !existe {
		return User{}, ErrUnauthorized
	}
	return s.users[username], nil
}