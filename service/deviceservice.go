package service

import (
	"encoding/base64"

	"github.com/chuckiihub/signing-service/crypto"
	"github.com/chuckiihub/signing-service/domain"
	apperrors "github.com/chuckiihub/signing-service/errors"
	"github.com/chuckiihub/signing-service/persistence"
)

// Persistence layer
type DeviceServiceImplementation struct {
	persistence persistence.DevicePersistance
	pageSize    int
}

// Creates a new device, assigns the key pair and saves it to storage
func (deviceService *DeviceServiceImplementation) Create(uuid string, algorithm crypto.SignatureAlgorithm, label string) (*domain.Device, error) {
	device := &domain.Device{
		UUID:          uuid,
		Algorithm:     algorithm,
		Label:         label,
		LastSignature: base64.StdEncoding.EncodeToString([]byte(uuid)),
	}

	crypto, err := crypto.NewCrypto(device.Algorithm)
	if err != nil {
		return nil, apperrors.WrapError(err, apperrors.InternalError)
	}

	keyPair, err := crypto.GenerateKeyPair()
	if err != nil {
		return nil, apperrors.WrapError(err, apperrors.InternalError)
	}

	publicKey, privateKey, err := crypto.Marshal(keyPair)
	if err != nil {
		return nil, apperrors.WrapError(err, apperrors.InternalError)
	}

	device.PrivateKey = privateKey
	device.PublicKey = publicKey
	device, err = deviceService.persistence.Save(device)
	if err != nil {
		return nil, apperrors.WrapError(err, apperrors.InternalError)
	}

	return device, nil
}

// Gets a device from storage and retrieves it.
func (deviceService *DeviceServiceImplementation) Get(uuid string) (*domain.Device, error) {
	device, err := deviceService.persistence.FindByUUID(uuid)

	return device, err
}

// Get all devices from storage and retrieves them.
func (deviceService *DeviceServiceImplementation) List(page int) ([]domain.Device, error) {
	devices, err := deviceService.persistence.List(page, deviceService.pageSize)

	if err != nil {
		return nil, apperrors.WrapError(err, apperrors.InternalError)
	}

	return devices, nil
}

func (deviceService *DeviceServiceImplementation) CheckHealth() domain.ServiceHealth {
	health := domain.ServiceHealth{PersistenceLayer: make(map[string]domain.PersistenceHealth)}

	dbHealth := deviceService.persistence.CheckHealth()
	health.PersistenceLayer["device"] = dbHealth
	health.Status = dbHealth.Status

	return health
}
