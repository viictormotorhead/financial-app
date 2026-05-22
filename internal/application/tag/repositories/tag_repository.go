package repositories

import (
	"context"

	"github.com/viictormotorhead/financial-app/internal/infra/repositories/tag/entities"
)

type TagRepositoryIF interface {
	Save(ctx context.Context, tag entities.TagEntity) (entities.TagEntity, error)
	ListAll(ctx context.Context) ([]entities.TagEntity, error)
	ResolveNames(ctx context.Context, names []string) ([]string, error)
}
