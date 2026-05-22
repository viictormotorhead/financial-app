package use_cases

import (
	"context"
	"fmt"

	"github.com/viictormotorhead/financial-app/internal/application/tag/dto/outputs"
	"github.com/viictormotorhead/financial-app/internal/application/tag/repositories"
)

type ListTagsUseCaseIF interface {
	List(ctx context.Context) (outputs.ListTagsOutputDTO, error)
}

type ListTagsUseCaseImpl struct {
	repository repositories.TagRepositoryIF
}

func NewListTagsUseCase(repository repositories.TagRepositoryIF) ListTagsUseCaseIF {
	return &ListTagsUseCaseImpl{repository: repository}
}

func (u *ListTagsUseCaseImpl) List(ctx context.Context) (outputs.ListTagsOutputDTO, error) {
	tags, err := u.repository.ListAll(ctx)
	if err != nil {
		return outputs.ListTagsOutputDTO{}, fmt.Errorf("list tags: %w", err)
	}

	items := make([]outputs.TagOutputDTO, 0, len(tags))
	for _, tag := range tags {
		items = append(items, outputs.TagOutputDTO{
			ID:          tag.ID,
			Name:        tag.Name,
			Description: tag.Description,
		})
	}

	return outputs.ListTagsOutputDTO{Tags: items}, nil
}
