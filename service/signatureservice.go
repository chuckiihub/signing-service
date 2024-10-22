package service

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/chuckiihub/signing-service/crypto"
	"github.com/chuckiihub/signing-service/domain"
	apperrors "github.com/chuckiihub/signing-service/errors"
	"github.com/chuckiihub/signing-service/persistence"
)

type SignatureServiceImplementation struct {
	devicePersistence    persistence.DevicePersistance
	signaturePersistence persistence.SignaturePersistance
	lockService          LockService
	pageSize             int
}

func (signingService *SignatureServiceImplementation) fetchDeviceOrReturnNotFound(deviceId string) (*domain.Device, error) {
	device, err := signingService.devicePersistence.FindByUUID(deviceId)
	if err != nil {
		return nil, apperrors.WrapError(err, apperrors.InternalError)
	}

	if device == nil {
		return nil, apperrors.WrapError(errors.New("device not found"), apperrors.NotFound)
	}

	return device, nil
}

func (signingService *SignatureServiceImplementation) Sign(deviceId string, dataToBeSigned string) (*domain.Signature, error) {
	if deviceId == "" {
		return nil, apperrors.WrapError(errors.New("deviceId is required"), apperrors.BadRequest)
	}

	signingService.lockService.Lock(deviceId)
	defer signingService.lockService.Unlock(deviceId)

	device, err := signingService.fetchDeviceOrReturnNotFound(deviceId)
	if err != nil {
		return nil, apperrors.WrapError(err, apperrors.InternalError)
	}
	originalDeviceCopy := *device
	device.SignatureCounter++

	dataToBeSigned = signingService.preSignEncoding(*device, dataToBeSigned)
	signature, err := signingService.signData(device, dataToBeSigned)
	if err != nil {
		return nil, err
	}

	signatureDTO := &domain.Signature{
		DeviceUUID: device.UUID,
		SignedData: dataToBeSigned,
		Signature:  base64.StdEncoding.EncodeToString(signature),
	}

	device.LastSignature = signatureDTO.Signature
	if _, err = signingService.devicePersistence.Save(device); err != nil {
		return nil, apperrors.WrapError(err, apperrors.InternalError)
	}

	signatureDTO, err = signingService.signaturePersistence.Save(signatureDTO)
	if err != nil {
		// rollbacks the device to its previous state as the second save failed.
		// ignores the error as it's not critical.
		signingService.devicePersistence.Save(&originalDeviceCopy)
		return nil, apperrors.WrapError(err, apperrors.InternalError)
	}

	return signatureDTO, nil
}

// This method should be called ALWAYS locking the device for writing using the
// LockingService. This protects the field SignatureCounter and LastSignature
// while signing requests.
func (signingService *SignatureServiceImplementation) signData(device *domain.Device, dataToBeSigned string) ([]byte, error) {

	crypto, err := crypto.NewCrypto(device.Algorithm)
	if err != nil {
		return nil, apperrors.WrapError(err, apperrors.InternalError)
	}

	signature, err := crypto.Sign([]byte(dataToBeSigned), []byte(device.PrivateKey))
	if err != nil {
		return nil, apperrors.WrapError(err, apperrors.InternalError)
	}

	return signature, nil
}

func (signingService *SignatureServiceImplementation) preSignEncoding(device domain.Device, data string) string {
	lastSignature := device.LastSignature
	if device.LastSignature == "" {
		lastSignature = device.UUID
	}

	return fmt.Sprintf("%d_%s_%s", device.SignatureCounter, data, base64.StdEncoding.EncodeToString([]byte(lastSignature)))
}

// TODO let's say that the person wants to see if the signature is valid or the data...
// Will they look up by the signature? By data? By Hash?
func (signingService *SignatureServiceImplementation) Get(signature string) (*domain.Signature, error) {
	signatureObject, err := signingService.signaturePersistence.FindBySignature(signature)
	if err != nil {
		return nil, apperrors.WrapError(err, apperrors.InternalError)
	}

	return signatureObject, nil
}

func (signingService *SignatureServiceImplementation) List(page int) ([]domain.Signature, error) {
	signatures, err := signingService.signaturePersistence.List(page, signingService.pageSize)

	if err != nil {
		return nil, apperrors.WrapError(err, apperrors.InternalError)
	}

	return signatures, nil
}

func (signingService *SignatureServiceImplementation) CheckHealth() domain.ServiceHealth {
	health := domain.ServiceHealth{Status: domain.HealthStatusPass, PersistenceLayer: map[string]domain.PersistenceHealth{}}
	deviceDbHealth := signingService.devicePersistence.CheckHealth()
	signatureDbHealth := signingService.signaturePersistence.CheckHealth()

	health.PersistenceLayer["device"] = deviceDbHealth
	health.PersistenceLayer["signature"] = signatureDbHealth

	if deviceDbHealth.Status == domain.HealthStatusPass && signatureDbHealth.Status == domain.HealthStatusPass {
		health.Status = domain.HealthStatusPass
	}

	return health
}

func (signingService *SignatureServiceImplementation) Verify(deviceId string, dataToBeSigned string, signature string) (bool, error) {
	device, err := signingService.fetchDeviceOrReturnNotFound(deviceId)
	if err != nil {
		return false, err
	}

	crypto, err := crypto.NewCrypto(device.Algorithm)
	if err != nil {
		return false, apperrors.WrapError(err, apperrors.InternalError)
	}

	// Decode the signature from base64 before verification
	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, apperrors.WrapError(err, apperrors.BadRequest)
	}

	// Verify using the decoded signature
	return crypto.Verify([]byte(dataToBeSigned), decodedSignature, device.PrivateKey)
}
