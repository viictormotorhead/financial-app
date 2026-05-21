package gormrepo

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	apptag "github.com/viictormotorhead/financial-app/internal/application/tag"
	"github.com/viictormotorhead/financial-app/internal/application/tag/repositories"
	"github.com/viictormotorhead/financial-app/internal/infra/repositories/tag/entities"
)

type tagWriteRepository struct {
	db *gorm.DB
}

func NewTagWriteRepository(db *gorm.DB) repositories.TagWriteRepositoryIF {
	return &tagWriteRepository{db: db}
}

func (r *tagWriteRepository) Save(ctx context.Context, tag entities.TagEntity) (entities.TagEntity, error) {
	model := entities.Tag{
		Name:        tag.Name,
		Description: tag.Description,
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		if isUniqueViolation(err) {
			return entities.TagEntity{}, apptag.ErrTagNameAlreadyExists
		}
		return entities.TagEntity{}, fmt.Errorf("insert tag: %w", err)
	}

	return entities.TagEntity{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
	}, nil
}

func isUniqueViolation(err error) bool {
	return errors.Is(err, gorm.ErrDuplicatedKey)
}
