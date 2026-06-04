package repositories

import (
	"context"

	"github.com/viictormotorhead/financial-app/internal/infra/repositories/tag/entities"
)

type TagRepositoryIF interface {
	Save(ctx context.Context, tag entities.TagEntity) (entities.TagEntity, error)
	ListAll(ctx context.Context, userID string) ([]entities.TagEntity, error)
	ResolveNames(ctx context.Context, userID string, names []string) ([]string, error)
}
