package gormrepo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	appauth "github.com/viictormotorhead/financial-app/internal/application/auth"
	appuser "github.com/viictormotorhead/financial-app/internal/application/user"
	"github.com/viictormotorhead/financial-app/internal/application/user/repositories"
	"github.com/viictormotorhead/financial-app/internal/infra/repositories/user/entities"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserRepositoryIF {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user entities.UserEntity) (entities.UserEntity, error) {
	now := time.Now().UTC()
	model := entities.User{
		UserID:       user.UserID,
		Name:         user.Name,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		CreatedAt:    now,
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return entities.UserEntity{}, appuser.ErrUsernameAlreadyExists
		}
		return entities.UserEntity{}, fmt.Errorf("insert user: %w", err)
	}

	return entities.UserEntity{
		UserID:   model.UserID,
		Name:     model.Name,
		Username: model.Username,
		CreatedAt: model.CreatedAt,
	}, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (entities.UserEntity, error) {
	var model entities.User
	err := r.db.WithContext(ctx).
		Where("LOWER(username) = ?", strings.ToLower(strings.TrimSpace(username))).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.UserEntity{}, appauth.ErrInvalidCredentials
		}
		return entities.UserEntity{}, fmt.Errorf("find user: %w", err)
	}

	return entities.UserEntity{
		UserID:       model.UserID,
		Name:         model.Name,
		Username:     model.Username,
		PasswordHash: model.PasswordHash,
		CreatedAt:    model.CreatedAt,
	}, nil
}
