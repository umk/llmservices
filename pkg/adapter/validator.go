package adapter

import (
	"github.com/go-playground/validator/v10"
	validatorutil "github.com/umk/llmservices/internal/validator"
)

func InitValidator(val *validator.Validate) {
	val.RegisterStructValidation(
		validatorutil.ValidateUnion,
		Message{},
		ContentPart{},
		ResponseFormat{},
	)
}
