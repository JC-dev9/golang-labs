package main

import (
	"database/sql"
	"log"
	"net/http"

	"desafio4/handlers"
	"desafio4/repository"
	"desafio4/services"

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
	
	// Rotas de Registo
	app.HandleFunc("GET /user/register", h.RegisterGET)
	app.HandleFunc("POST /user/register", h.RegisterPOST)

	// Rotas de Login
	app.HandleFunc("GET /user/login", h.LoginGET)
	app.HandleFunc("POST /user/login", h.LoginPOST)

	// Rotas Protegidas
	app.HandleFunc("GET /user/profile", h.Profile)
	app.HandleFunc("POST /user/logout", h.Logout)

	log.Println("Servidor web a correr em http://localhost:8080")
	if err := http.ListenAndServe(":8080", app); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}