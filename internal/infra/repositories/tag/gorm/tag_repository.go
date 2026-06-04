package gormrepo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	apptag "github.com/viictormotorhead/financial-app/internal/application/tag"
	"github.com/viictormotorhead/financial-app/internal/application/tag/repositories"
	"github.com/viictormotorhead/financial-app/internal/infra/repositories/tag/entities"
	"github.com/viictormotorhead/financial-app/internal/infra/repositories/scopes"
)

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) repositories.TagRepositoryIF {
	return &tagRepository{db: db}
}

func (r *tagRepository) Save(ctx context.Context, tag entities.TagEntity) (entities.TagEntity, error) {
	model := entities.Tag{
		UserID:      tag.UserID,
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
		UserID:      model.UserID,
		Name:        model.Name,
		Description: model.Description,
	}, nil
}

func (r *tagRepository) ListAll(ctx context.Context, userID string) ([]entities.TagEntity, error) {
	var tags []entities.Tag
	if err := scopes.UserID(r.db.WithContext(ctx), userID).Order("name ASC").Find(&tags).Error; err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}

	result := make([]entities.TagEntity, 0, len(tags))
	for _, tag := range tags {
		result = append(result, entities.TagEntity{
			ID:          tag.ID,
			UserID:      tag.UserID,
			Name:        tag.Name,
			Description: tag.Description,
		})
	}

	return result, nil
}

func (r *tagRepository) ResolveNames(ctx context.Context, userID string, names []string) ([]string, error) {
	if len(names) == 0 {
		return []string{}, nil
	}

	keys := make([]string, 0, len(names))
	seen := make(map[string]struct{}, len(names))
	for _, name := range names {
		key := tagNameKey(name)
		if key == "" {
			continue
		}
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		keys = append(keys, key)
	}

	if len(keys) == 0 {
		return []string{}, nil
	}

	var tags []entities.Tag
	if err := scopes.UserID(r.db.WithContext(ctx), userID).
		Where("LOWER(name) IN ?", keys).
		Find(&tags).Error; err != nil {
		return nil, fmt.Errorf("find tags: %w", err)
	}

	canonicalByKey := make(map[string]string, len(tags))
	for _, tag := range tags {
		canonicalByKey[tagNameKey(tag.Name)] = tag.Name
	}

	resolved := make([]string, 0, len(names))
	missing := make([]string, 0)
	seenResolved := make(map[string]struct{}, len(names))

	for _, name := range names {
		key := tagNameKey(name)
		if key == "" {
			continue
		}
		if _, exists := seenResolved[key]; exists {
			continue
		}
		seenResolved[key] = struct{}{}

		canonical, ok := canonicalByKey[key]
		if !ok {
			missing = append(missing, strings.TrimSpace(name))
			continue
		}
		resolved = append(resolved, canonical)
	}

	if len(missing) > 0 {
		return nil, apptag.NewTagsNotFoundError(missing)
	}

	return resolved, nil
}

func tagNameKey(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func isUniqueViolation(err error) bool {
	return errors.Is(err, gorm.ErrDuplicatedKey)
}
