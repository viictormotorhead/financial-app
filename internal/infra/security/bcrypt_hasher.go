package security

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	appauth "github.com/viictormotorhead/financial-app/internal/application/auth"
	"github.com/viictormotorhead/financial-app/internal/application/user/services"
)

type bcryptHasher struct {
	cost int
}

func NewBcryptHasher(cost int) services.PasswordHasher {
	if cost < bcrypt.MinCost {
		cost = bcrypt.DefaultCost
	}
	return &bcryptHasher{cost: cost}
}

func (h *bcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(hash), nil
}

func (h *bcryptHasher) Verify(hashedPassword, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err == nil {
		return nil
	}
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return appauth.ErrInvalidCredentials
	}
	return fmt.Errorf("verify password: %w", err)
}
