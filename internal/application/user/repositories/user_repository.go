package repositories

import (
	"context"

	"github.com/viictormotorhead/financial-app/internal/infra/repositories/user/entities"
)

type UserRepositoryIF interface {
	Create(ctx context.Context, user entities.UserEntity) (entities.UserEntity, error)
	FindByUsername(ctx context.Context, username string) (entities.UserEntity, error)
}
