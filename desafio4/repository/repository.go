package repository

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrSessionNotFound = errors.New("session not found")
)

// Entidades Canónicas 
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // O traço é para não expor a password em JSON
	CreatedAt time.Time `json:"created_at"`
}

type Session struct {
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Interfaces
type UserRepository interface {
	Create(user User) error
	FindByUsername(username string) (User, error)
	FindByID(id string) (User, error)
	FindAll() ([]User, error)
}

type SessionRepository interface {
	Create(session Session) error
	FindByToken(token string) (Session, error)
	DeleteByToken(token string) error
}