package use_cases

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	apptag "github.com/viictormotorhead/financial-app/internal/application/tag"
	tagrepositories "github.com/viictormotorhead/financial-app/internal/application/tag/repositories"
	"github.com/viictormotorhead/financial-app/internal/application/investment/dto/outputs"
	"github.com/viictormotorhead/financial-app/internal/application/investment/repositories"
	"github.com/viictormotorhead/financial-app/internal/infra/repositories/investment/entities"
)

type CreateInvestmentCommand struct {
	Name    string
	Balance float64
	Tags    []string
}

type CreateInvestmentUseCaseIF interface {
	Create(ctx context.Context, cmd CreateInvestmentCommand) (outputs.CreateInvestmentOutputDTO, error)
}

type CreateInvestmentUseCaseImpl struct {
	repository    repositories.InvestmentWriteRepositoryIF
	tagRepository tagrepositories.TagRepositoryIF
}

func NewCreateInvestmentUseCase(
	repository repositories.InvestmentWriteRepositoryIF,
	tagRepository tagrepositories.TagRepositoryIF,
) CreateInvestmentUseCaseIF {
	return &CreateInvestmentUseCaseImpl{
		repository:    repository,
		tagRepository: tagRepository,
	}
}

func (u *CreateInvestmentUseCaseImpl) Create(ctx context.Context, cmd CreateInvestmentCommand) (outputs.CreateInvestmentOutputDTO, error) {
	tags := normalizeTags(cmd.Tags)

	if len(tags) > 0 {
		resolved, err := u.tagRepository.ResolveNames(ctx, tags)
		if err != nil {
			var tagsNotFound *apptag.TagsNotFoundError
			if errors.As(err, &tagsNotFound) {
				return outputs.CreateInvestmentOutputDTO{}, err
			}
			return outputs.CreateInvestmentOutputDTO{}, fmt.Errorf("resolve tags: %w", err)
		}
		tags = resolved
	}

	entity := entities.InvestmentEntity{
		Name:           cmd.Name,
		Balance:        cmd.Balance,
		InitialBalance: cmd.Balance,
		Tags:           tags,
		CreatedAt:      time.Now().UTC(),
	}

	saved, err := u.repository.Save(ctx, entity)
	if err != nil {
		return outputs.CreateInvestmentOutputDTO{}, fmt.Errorf("create investment: %w", err)
	}

	return outputs.CreateInvestmentOutputDTO{
		ID:             saved.ID,
		Name:           saved.Name,
		Balance:        saved.Balance,
		InitialBalance: saved.InitialBalance,
		Tags:           saved.Tags,
		CreatedAt:      saved.CreatedAt,
	}, nil
}

func normalizeTags(tags []string) []string {
	if len(tags) == 0 {
		return []string{}
	}

	seen := make(map[string]struct{}, len(tags))
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		if trimmed == "" {
			continue
		}
		key := strings.ToLower(trimmed)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, trimmed)
	}

	return result
}
