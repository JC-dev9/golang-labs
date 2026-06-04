package repository_test

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"desafio3/repository"
	_ "modernc.org/sqlite" // O driver
)

// Helper para criar uma BD limpa antes de cada teste
func setupTestDB(t *testing.T) (*sql.DB, *repository.SQLiteUserRepository) {
	t.Helper()

	// ficheiro temporário que o Go apaga no fim do teste
	dbPath := filepath.Join(t.TempDir(), "test.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("erro ao abrir bd de teste: %v", err)
	}

	query := `
	CREATE TABLE users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE,
		password TEXT,
		created_at DATETIME
	)`
	if _, err := db.Exec(query); err != nil {
		t.Fatalf("erro ao criar tabela de teste: %v", err)
	}

	repo := repository.NewSQLiteUserRepository(db)
	
	return db, repo
}

func TestSQLiteUserRepository(t *testing.T) {
	db, repo := setupTestDB(t)
	
	// O defer garante que fechamos a ligação à BD quando o teste acabar
	defer db.Close()

	t.Run("consegue criar e buscar um utilizador na BD real", func(t *testing.T) {

		novoUser := repository.User{
			ID:        "sql-123",
			Username:  "juan_sqlite",
			Password:  "hash_seguro",
			CreatedAt: time.Now(),
		}

		err := repo.Create(novoUser)
		if err != nil {
			t.Fatalf("não esperava erro ao criar user no sqlite: %v", err)
		}

		userRetornado, err := repo.FindByUsername("juan_sqlite")
		if err != nil {
			t.Fatalf("erro ao buscar user por username: %v", err)
		}
		
		if userRetornado.ID != "sql-123" {
			t.Errorf("esperava ID sql-123, recebi %s", userRetornado.ID)
		}

		userPorID, err := repo.FindByID("sql-123")
		if err != nil {
			t.Fatalf("erro ao buscar user por ID: %v", err)
		}

		if userPorID.Username != "juan_sqlite" {
			t.Errorf("esperava Username juan_sqlite, recebi %s", userPorID.Username)
		}
	})
}