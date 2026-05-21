package request

import (
	"errors"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateInvestmentRequest struct {
	Name    string   `json:"name"`
	Balance float64  `json:"balance"`
	Tags    []string `json:"tags"`
}

func (r CreateInvestmentRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name,
			validation.Required.Error("name is required"),
			validation.Length(1, 255).Error("name must be between 1 and 255 characters"),
		),
		validation.Field(&r.Balance,
			validation.Required.Error("balance is required"),
			validation.Min(0.01).Error("balance must be greater than zero"),
		),
		validation.Field(&r.Tags,
			validation.Each(
				validation.Required.Error("tag cannot be empty"),
				validation.Length(1, 255).Error("each tag must be between 1 and 255 characters"),
			),
			validation.By(validateUniqueTags),
		),
	)
}

func validateUniqueTags(value interface{}) error {
	tags, ok := value.([]string)
	if !ok || len(tags) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		normalized := strings.TrimSpace(strings.ToLower(tag))
		if normalized == "" {
			continue
		}
		if _, exists := seen[normalized]; exists {
			return errors.New("tags must not contain duplicates")
		}
		seen[normalized] = struct{}{}
	}

	return nil
}
