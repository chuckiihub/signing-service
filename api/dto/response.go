package dto

import (
	"github.com/chuckiihub/signing-service/domain"
)

// Represents the server's response to the client's request to create a new device.
type DeviceResponse struct {
	Id         string `json:"uuid"`
	Label      string `json:"label"`
	Algorithm  string `json:"algorithm"`
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

type SignatureResponse struct {
	DeviceId   string `json:"deviceId"`
	SignedData string `json:"signedData"`
	Signature  string `json:"signature"`
}

func NewSignatureResponseFromSignature(signature *domain.Signature) *SignatureResponse {
	return &SignatureResponse{
		DeviceId:   signature.DeviceUUID,
		SignedData: signature.SignedData,
		Signature:  signature.Signature,
	}
}

func NewSignatureResponse(signature *domain.Signature) SignatureResponse {
	return SignatureResponse{
		DeviceId:   signature.DeviceUUID,
		SignedData: signature.SignedData,
		Signature:  signature.Signature,
	}
}

func NewDeviceResponse(device *domain.Device) DeviceResponse {
	publicKeyPEM := string(device.PublicKey)

	return DeviceResponse{
		Id:         device.UUID,
		Label:      device.Label,
		Algorithm:  device.Algorithm.String(),
		PublicKey:  publicKeyPEM,
		PrivateKey: string(device.PrivateKey),
	}
}
