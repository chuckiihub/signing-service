package validation

import (
	"slices"

	"github.com/chuckiihub/signing-service/api/dto"
	"github.com/chuckiihub/signing-service/crypto"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type RequestValidator struct {
	validationService *validator.Validate
}

func (validationService *RequestValidator) GetValidationFailureErrors(err error) []string {
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

func (validationService *RequestValidator) Validate(request interface{}) error {
	return validationService.validationService.Struct(request)
}

func (validationService *RequestValidator) ValidateDeviceCreationRequest(request dto.DeviceCreationRequest) error {
	return validationService.validationService.Struct(request)
}

func (validationService *RequestValidator) ValidateSignatureCreationRequest(request dto.SignatureCreateRequest) error {
	return validationService.validationService.Struct(request)
}

func (validationService *RequestValidator) ValidateSignatureVerificationRequest(request dto.SignatureVerifyRequest) error {
	return validationService.validationService.Struct(request)
}

func NewRequestValidator() RequestValidator {
	validate = validator.New()
	validate.RegisterValidation("supported-encryption", validateSignatureAlgorithm)

	return RequestValidator{validationService: validate}
}

// Validates automatically that the string is one of SignatureAlgorithm enum
func validateSignatureAlgorithm(fieldLevel validator.FieldLevel) bool {
	algorithm := fieldLevel.Field().String()

	return slices.IndexFunc(crypto.GetSupportedAlgorithms(), func(alg string) bool { return alg == algorithm }) >= 0
}
