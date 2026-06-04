package repository_test

import (
	"desafio4/repository"
	"testing"
)

func TestMemoryUserRepository(t *testing.T) {
	repo := repository.NewMemoryUserRepository()

	t.Run("consegue criar e buscar um utilizador pelo username", func(t *testing.T) {
		
		// AAA - Arrange, Act, Assert

		novoUser := repository.User{
			ID:       "12345678",
			Username: "juan_loza",
			Password: "hashedpassword123",
		}
		
		err := repo.Create(novoUser)

		if err != nil {
			t.Fatalf("não esperava erro ao criar user, mas recebi: %v", err)
		}

		userRetornado, err := repo.FindByUsername("juan_loza")

		if err != nil {
			t.Fatalf("erro ao buscar user: %v", err)
		}
		if userRetornado.ID != "12345678" {
			t.Errorf("esperava 12345678, recebi %s", userRetornado.ID)
		}
	})
}