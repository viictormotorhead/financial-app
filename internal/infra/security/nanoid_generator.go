package security

import (
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"

	"github.com/viictormotorhead/financial-app/internal/application/user/services"
)

type nanoidGenerator struct{}

func NewNanoidGenerator() services.IDGenerator {
	return &nanoidGenerator{}
}

func (g *nanoidGenerator) New() (string, error) {
	id, err := gonanoid.New()
	if err != nil {
		return "", fmt.Errorf("generate user id: %w", err)
	}
	return id, nil
}
