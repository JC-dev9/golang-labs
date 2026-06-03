package repository

import "errors"

type MemoryUserRepository struct {
	users map[string]User
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: make(map[string]User),
	}
}

func (m *MemoryUserRepository) Create(user User) error {

	if _, exists := m.users[user.Username]; exists {
		return errors.New("user já existe")
	}
	m.users[user.Username] = user	
	return nil
}

func (m *MemoryUserRepository) FindByUsername(username string) (User, error) {

	if _, exists := m.users[username]; !exists {
		return User{}, ErrUserNotFound
	}
	return m.users[username], nil
}

func (m *MemoryUserRepository) FindByID(id string) (User, error) {
	// Como o mapa é indexado por username, temos de o percorrer:
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return User{}, ErrUserNotFound
}

func (m *MemoryUserRepository) FindAll() ([]User, error) {
	var users []User
	
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}


type MemorySessionRepository struct {
	sessions map[string]Session
}

func NewMemorySessionRepository() *MemorySessionRepository {
	return &MemorySessionRepository{
		sessions: make(map[string]Session),
	}
}

func (m *MemorySessionRepository) Create(session Session) error {

	if _, exists := m.sessions[session.Token]; exists {
		return errors.New("sessão já existe")
	}
	m.sessions[session.Token] = session
	return nil 
}

func (m *MemorySessionRepository) FindByToken(token string) (Session, error) {
	
	if _, existe := m.sessions[token]; !existe {
		return Session{}, ErrSessionNotFound
	}

	return m.sessions[token], nil 
}

func (m *MemorySessionRepository) DeleteByToken(token string) error {

	delete(m.sessions, token)
	return nil 
}

