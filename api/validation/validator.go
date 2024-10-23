package validation

import (
	"slices"

	"github.com/chuckiihub/signing-service/crypto"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type RequestValidator struct {
	validator *validator.Validate
}

func (validation *RequestValidator) GetValidationFailureErrors(err error) []string {
	var errors []string
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, validationErr := range validationErrors {
			errors = append(errors, validationErr.Field()+" failed the following validation: "+validationErr.Tag())

		}
	} else {
		errors = append(errors, err.Error())
	}

	return errors
}

func (validation *RequestValidator) Validate(request interface{}) error {
	return validation.validator.Struct(request)
}

func NewRequestValidator() RequestValidator {
	validate = validator.New()
	validate.RegisterValidation("supported-encryption", validateSignatureAlgorithm)

	return RequestValidator{validator: validate}
}

// Validates automatically that the string is one of SignatureAlgorithm enum
func validateSignatureAlgorithm(fieldLevel validator.FieldLevel) bool {
	algorithm := fieldLevel.Field().String()

	return slices.IndexFunc(crypto.GetSupportedAlgorithms(), func(alg string) bool { return alg == algorithm }) >= 0
}
