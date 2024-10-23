package dto

import (
	"testing"

	"github.com/chuckiihub/signing-service/crypto"
	"github.com/chuckiihub/signing-service/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewSignatureResponseFromSignature(t *testing.T) {
	signature := &domain.Signature{
		DeviceUUID: "test-device-uuid",
		SignedData: "test-signed-data",
		Signature:  "test-signature",
	}

	response := NewSignatureResponseFromSignature(signature)

	assert.Equal(t, signature.DeviceUUID, response.DeviceId)
	assert.Equal(t, signature.SignedData, response.SignedData)
	assert.Equal(t, signature.Signature, response.Signature)
}

func TestNewSignatureResponse(t *testing.T) {
	signature := &domain.Signature{
		DeviceUUID: "test-device-uuid",
		SignedData: "test-signed-data",
		Signature:  "test-signature",
	}

	response := NewSignatureResponse(signature)

	assert.Equal(t, signature.DeviceUUID, response.DeviceId)
	assert.Equal(t, signature.SignedData, response.SignedData)
	assert.Equal(t, signature.Signature, response.Signature)
}

func TestNewDeviceResponse(t *testing.T) {
	device := &domain.Device{
		UUID:       "test-uuid",
		Label:      "test-label",
		Algorithm:  crypto.SignatureAlgorithmRSA,
		PublicKey:  []byte("test-public-key"),
		PrivateKey: []byte("test-private-key"),
	}

	response := NewDeviceResponse(device)

	assert.Equal(t, device.UUID, response.Id)
	assert.Equal(t, device.Label, response.Label)
	assert.Equal(t, device.Algorithm.String(), response.Algorithm)
	assert.Equal(t, string(device.PublicKey), response.PublicKey)
	assert.Equal(t, string(device.PrivateKey), response.PrivateKey)
}
