package use_cases

import (
	"context"
	"errors"
	"fmt"
	"strings"

	appauth "github.com/viictormotorhead/financial-app/internal/application/auth"
	appuser "github.com/viictormotorhead/financial-app/internal/application/user"
	"github.com/viictormotorhead/financial-app/internal/application/user/dto/outputs"
	"github.com/viictormotorhead/financial-app/internal/application/user/repositories"
	"github.com/viictormotorhead/financial-app/internal/application/user/services"
	userentities "github.com/viictormotorhead/financial-app/internal/infra/repositories/user/entities"
)

type RegisterUserCommand struct {
	Name     string
	Username string
	Password string
}

type RegisterUserUseCaseIF interface {
	Register(ctx context.Context, cmd RegisterUserCommand) (outputs.AuthOutputDTO, error)
}

type RegisterUserUseCaseImpl struct {
	repository     repositories.UserRepositoryIF
	passwordHasher services.PasswordHasher
	idGenerator    services.IDGenerator
	tokenIssuer    appauth.TokenIssuer
}

func NewRegisterUserUseCase(
	repository repositories.UserRepositoryIF,
	passwordHasher services.PasswordHasher,
	idGenerator services.IDGenerator,
	tokenIssuer appauth.TokenIssuer,
) RegisterUserUseCaseIF {
	return &RegisterUserUseCaseImpl{
		repository:     repository,
		passwordHasher: passwordHasher,
		idGenerator:    idGenerator,
		tokenIssuer:    tokenIssuer,
	}
}

func (u *RegisterUserUseCaseImpl) Register(ctx context.Context, cmd RegisterUserCommand) (outputs.AuthOutputDTO, error) {
	userID, err := u.idGenerator.New()
	if err != nil {
		return outputs.AuthOutputDTO{}, fmt.Errorf("generate user id: %w", err)
	}

	passwordHash, err := u.passwordHasher.Hash(cmd.Password)
	if err != nil {
		return outputs.AuthOutputDTO{}, fmt.Errorf("hash password: %w", err)
	}

	saved, err := u.repository.Create(ctx, userentities.UserEntity{
		UserID:       userID,
		Name:         strings.TrimSpace(cmd.Name),
		Username:     normalizeUsername(cmd.Username),
		PasswordHash: passwordHash,
	})
	if err != nil {
		if errors.Is(err, appuser.ErrUsernameAlreadyExists) {
			return outputs.AuthOutputDTO{}, err
		}
		return outputs.AuthOutputDTO{}, fmt.Errorf("create user: %w", err)
	}

	token, err := u.tokenIssuer.Issue(ctx, appauth.TokenClaims{
		UserID:   saved.UserID,
		Username: saved.Username,
	})
	if err != nil {
		return outputs.AuthOutputDTO{}, fmt.Errorf("issue token: %w", err)
	}

	return outputs.AuthOutputDTO{
		Token: token,
		User: outputs.AuthUserOutputDTO{
			UserID:   saved.UserID,
			Name:     saved.Name,
			Username: saved.Username,
		},
	}, nil
}

func normalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}
