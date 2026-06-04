package use_cases

import (
	"context"
	"errors"
	"fmt"

	appauth "github.com/viictormotorhead/financial-app/internal/application/auth"
	"github.com/viictormotorhead/financial-app/internal/application/user/dto/outputs"
	"github.com/viictormotorhead/financial-app/internal/application/user/repositories"
	"github.com/viictormotorhead/financial-app/internal/application/user/services"
)

type LoginUserCommand struct {
	Username string
	Password string
}

type LoginUserUseCaseIF interface {
	Login(ctx context.Context, cmd LoginUserCommand) (outputs.AuthOutputDTO, error)
}

type LoginUserUseCaseImpl struct {
	repository     repositories.UserRepositoryIF
	passwordHasher services.PasswordHasher
	tokenIssuer    appauth.TokenIssuer
}

func NewLoginUserUseCase(
	repository repositories.UserRepositoryIF,
	passwordHasher services.PasswordHasher,
	tokenIssuer appauth.TokenIssuer,
) LoginUserUseCaseIF {
	return &LoginUserUseCaseImpl{
		repository:     repository,
		passwordHasher: passwordHasher,
		tokenIssuer:    tokenIssuer,
	}
}

func (u *LoginUserUseCaseImpl) Login(ctx context.Context, cmd LoginUserCommand) (outputs.AuthOutputDTO, error) {
	user, err := u.repository.FindByUsername(ctx, cmd.Username)
	if err != nil {
		if errors.Is(err, appauth.ErrInvalidCredentials) {
			return outputs.AuthOutputDTO{}, err
		}
		return outputs.AuthOutputDTO{}, fmt.Errorf("find user: %w", err)
	}

	if err := u.passwordHasher.Verify(user.PasswordHash, cmd.Password); err != nil {
		if errors.Is(err, appauth.ErrInvalidCredentials) {
			return outputs.AuthOutputDTO{}, err
		}
		return outputs.AuthOutputDTO{}, fmt.Errorf("verify password: %w", err)
	}

	token, err := u.tokenIssuer.Issue(ctx, appauth.TokenClaims{
		UserID:   user.UserID,
		Username: user.Username,
	})
	if err != nil {
		return outputs.AuthOutputDTO{}, fmt.Errorf("issue token: %w", err)
	}

	return outputs.AuthOutputDTO{
		Token: token,
		User: outputs.AuthUserOutputDTO{
			UserID:   user.UserID,
			Name:     user.Name,
			Username: user.Username,
		},
	}, nil
}
