package repository

import (
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

// repositório de utilizadores

type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

func (s *SQLiteUserRepository) Create(user User) error {
	query := `INSERT INTO users (id, username, password, created_at) VALUES (?, ?, ?, ?)`
	_, err := s.db.Exec(query, user.ID, user.Username, user.Password, user.CreatedAt)
	return err
}

func (s *SQLiteUserRepository) FindByUsername(username string) (User, error) {
	var u User
	query := `SELECT id, username, password, created_at FROM users WHERE username = ?`
	
	err := s.db.QueryRow(query, username).Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return u, nil
}

func (s *SQLiteUserRepository) FindByID(id string) (User, error) {
	var u User
	query := `SELECT id, username, password, created_at FROM users WHERE id = ?`
	
	err := s.db.QueryRow(query, id).Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return u, nil
}

func (s *SQLiteUserRepository) FindAll() ([]User, error) {

	query := `SELECT id, username, password, created_at FROM users`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User

		err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// repositório de sessões 

type SQLiteSessionRepository struct {
	db *sql.DB
}

func NewSQLiteSessionRepository(db *sql.DB) *SQLiteSessionRepository {
	return &SQLiteSessionRepository{db: db}
}

func (s *SQLiteSessionRepository) Create(session Session) error {
	query := `INSERT INTO sessions (token, user_id, created_at) VALUES (?, ?, ?)`
	_, err := s.db.Exec(query, session.Token, session.UserID, session.CreatedAt)
	return err
}

func (s *SQLiteSessionRepository) FindByToken(token string) (Session, error) {
	var sess Session
	query := `SELECT token, user_id, created_at FROM sessions WHERE token = ?`
	
	err := s.db.QueryRow(query, token).Scan(&sess.Token, &sess.UserID, &sess.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Session{}, ErrSessionNotFound
		}
		return Session{}, err
	}
	return sess, nil
}

func (s *SQLiteSessionRepository) DeleteByToken(token string) error {
	query := `DELETE FROM sessions WHERE token = ?`
	_, err := s.db.Exec(query, token)
	return err
}