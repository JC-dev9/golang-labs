package main

import (
	"database/sql"
	"log"
	"net/http"

	"desafio3/handlers"
	"desafio3/repository"
	"desafio3/services"

	_ "modernc.org/sqlite" // O driver
)

// Helper para manter o main() limpo
func createTables(db *sql.DB) {
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE,
		password TEXT,
		created_at DATETIME
	);`

	sessionsTable := `
	CREATE TABLE IF NOT EXISTS sessions (
		token TEXT PRIMARY KEY,
		user_id TEXT,
		created_at DATETIME
	);`

	if _, err := db.Exec(usersTable); err != nil {
		log.Fatalf("Erro ao criar tabela users: %v", err)
	}
	if _, err := db.Exec(sessionsTable); err != nil {
		log.Fatalf("Erro ao criar tabela sessions: %v", err)
	}
}

func main() {
	db, err := sql.Open("sqlite", "auth.db")
	if err != nil {
		log.Fatalf("Erro ao ligar à base de dados: %v", err)
	}
	defer db.Close()

	createTables(db)

	// Criar os repositórios reais que usam SQLite
	userRepo := repository.NewSQLiteUserRepository(db)
	sessionRepo := repository.NewSQLiteSessionRepository(db)
	
	authService := services.NewAuthService(userRepo, sessionRepo)
	
	h := handlers.NewHandlers(authService)

	app := http.NewServeMux()
	app.HandleFunc("/api/user/register", h.Register)
	app.HandleFunc("/api/user/login", h.Login)
	app.HandleFunc("/api/user/profile", h.Profile)

	log.Println("Servidor a correr na porta 8080 com base de dados SQLite")
	if err := http.ListenAndServe(":8080", app); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
