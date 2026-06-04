package use_cases

import (
	"context"
	"errors"
	"fmt"

	"github.com/viictormotorhead/financial-app/internal/application/auth"
	apptag "github.com/viictormotorhead/financial-app/internal/application/tag"
	"github.com/viictormotorhead/financial-app/internal/application/tag/dto/outputs"
	"github.com/viictormotorhead/financial-app/internal/application/tag/repositories"
	"github.com/viictormotorhead/financial-app/internal/infra/repositories/tag/entities"
)

type CreateTagCommand struct {
	Name        string
	Description string
}

type CreateTagUseCaseIF interface {
	Create(ctx context.Context, cmd CreateTagCommand) (outputs.CreateTagOutputDTO, error)
}

type CreateTagUseCaseImpl struct {
	repository repositories.TagRepositoryIF
}

func NewCreateTagUseCase(repository repositories.TagRepositoryIF) CreateTagUseCaseIF {
	return &CreateTagUseCaseImpl{repository: repository}
}

func (u *CreateTagUseCaseImpl) Create(ctx context.Context, cmd CreateTagCommand) (outputs.CreateTagOutputDTO, error) {
	userID, err := auth.RequireUserID(ctx)
	if err != nil {
		return outputs.CreateTagOutputDTO{}, err
	}

	entity := entities.TagEntity{
		UserID:      auth.UserIDPtr(userID),
		Name:        cmd.Name,
		Description: cmd.Description,
	}

	saved, err := u.repository.Save(ctx, entity)
	if err != nil {
		if errors.Is(err, apptag.ErrTagNameAlreadyExists) {
			return outputs.CreateTagOutputDTO{}, apptag.ErrTagNameAlreadyExists
		}
		return outputs.CreateTagOutputDTO{}, fmt.Errorf("create tag: %w", err)
	}

	return outputs.CreateTagOutputDTO{
		ID:          saved.ID,
		Name:        saved.Name,
		Description: saved.Description,
	}, nil
}
