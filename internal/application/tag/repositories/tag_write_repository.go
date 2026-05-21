package repositories

import (
	"context"

	"github.com/viictormotorhead/financial-app/internal/infra/repositories/tag/entities"
)

type TagWriteRepositoryIF interface {
	Save(ctx context.Context, tag entities.TagEntity) (entities.TagEntity, error)
}
