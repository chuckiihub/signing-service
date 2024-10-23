package dto

import (
	"errors"

	"github.com/chuckiihub/signing-service/crypto"
)

// Represents the client's request to create a new device.
type DeviceCreationRequest struct {
	Label     string `json:"label" validate:"required,min=1"`
	Algorithm string `json:"algorithm" validate:"required,min=1,supported-encryption"`
}

// Retrieves the string representation of the algorithm into the corresponding domain type.
func (request *DeviceCreationRequest) GetSignatureAlgorithm() (crypto.SignatureAlgorithm, error) {
	switch request.Algorithm {
	case "RSA":
		return crypto.SignatureAlgorithmRSA, nil
	case "ECC":
		return crypto.SignatureAlgorithmECC, nil
	default:
		return -1, errors.New("algorithm not supported")
	}
}

// Client request to sign new data
type SignatureCreateRequest struct {
	Data string `json:"data" validate:"required"`
}

// Client request to verified already signed data
type SignatureVerifyRequest struct {
	SignedData string `json:"signedData" validate:"required,min=1"`
	Signature  string `json:"signature" validate:"required,min=1"`
}
