package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"

	"github.com/chuckiihub/signing-service/crypto"
	"github.com/chuckiihub/signing-service/domain"
	apperrors "github.com/chuckiihub/signing-service/errors"
	"github.com/chuckiihub/signing-service/persistence"
	"github.com/google/uuid"
)

type SignatureServiceImplementation struct {
	devicePersistence    persistence.DevicePersistance
	signaturePersistence persistence.SignaturePersistance
	lockService          LockService
	pageSize             int
}

// Handy method to fetch a device and check errors.
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
	newSignatureUUID := uuid.NewString()

	// I check before the lock so we don't use the locking service in vain
	// in case of, for example, a DoS attack with non existing deviceIds.
	_, err := signingService.fetchDeviceOrReturnNotFound(deviceId)
	if err != nil {
		return nil, apperrors.WrapError(err, apperrors.InternalError)
	}

	signingService.lockService.Lock(deviceId)
	defer signingService.lockService.Unlock(deviceId)

	// after obtaining the lock, I need to refetch the device to ensure it wasn't changed meanwhile.
	device, err := signingService.fetchDeviceOrReturnNotFound(deviceId)
	if err != nil {
		return nil, apperrors.WrapError(err, apperrors.InternalError)
	}
	originalLastSignature := device.LastSignature
	originalSignatureCounter := device.SignatureCounter

	// no need to increment using atomic package as said in the requirements as it's protected by the lock
	device.SignatureCounter++

	dataToBeSigned = signingService.preSignEncoding(*device, dataToBeSigned)
	signature, err := signingService.signData(device, dataToBeSigned)
	if err != nil {
		slog.Warn("error while signing data", "error", err.Error())
		return nil, err
	}

	signatureDTO := &domain.Signature{
		UUID:       newSignatureUUID,
		DeviceUUID: device.UUID,
		SignedData: dataToBeSigned,
		Signature:  base64.StdEncoding.EncodeToString(signature),
	}

	device.LastSignature = signatureDTO.Signature
	if _, err = signingService.devicePersistence.Save(device); err != nil {
		// If saving the device fails, we discard the newly created signature.
		return nil, apperrors.WrapError(err, apperrors.InternalError)
	}

	signatureDTO, err = signingService.signaturePersistence.Save(signatureDTO)
	if err != nil {
		// If saving the signature fails, we rollbacks the device to its previous state.
		device.SignatureCounter = originalSignatureCounter
		device.LastSignature = originalLastSignature
		signingService.devicePersistence.Save(device)

		slog.Warn("error saving signature, trying to rollback device last signature and signatureCounter", "error", err.Error())
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
		slog.Warn("error while trying to sign Signer", "error", err.Error())
		return nil, apperrors.WrapError(err, apperrors.InternalError)
	}

	return signature, nil
}

// This method is used to concatenate the signatureCounter and the lastSignature to the data to be signed
// to create a unique signature for each device.
// If necessary, in the future, it could be moved to an Encoding package so support for other encodings is added.
func (signingService *SignatureServiceImplementation) preSignEncoding(device domain.Device, data string) string {
	lastSignature := device.LastSignature
	if device.LastSignature == "" {
		lastSignature = device.UUID
	}

	return fmt.Sprintf("%d_%s_%s", device.SignatureCounter, data, base64.StdEncoding.EncodeToString([]byte(lastSignature)))
}

func (signingService *SignatureServiceImplementation) Get(uuid string) (*domain.Signature, error) {
	signatureObject, err := signingService.signaturePersistence.FindByUUID(uuid)
	if err != nil {
		return nil, apperrors.WrapError(err, apperrors.InternalError)
	}

	return signatureObject, nil
}

func (signingService *SignatureServiceImplementation) List(page int) ([]domain.Signature, error) {
	if page < 1 {
		page = 1
	}

	signatures, err := signingService.signaturePersistence.List(page, signingService.pageSize)

	if err != nil {
		return nil, apperrors.WrapError(err, apperrors.InternalError)
	}

	return signatures, nil
}

// Check the health of the dependencies of this service.
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

// This method will receive the full signedData and signature and return a boolean
// indicating if the signature is valid for the given dataToBeSigned.
// If there's any error returned the signature is not valid.
func (signingService *SignatureServiceImplementation) Verify(deviceId string, dataToBeSigned string, signature string) (bool, error) {
	device, err := signingService.fetchDeviceOrReturnNotFound(deviceId)
	if err != nil {
		return false, err
	}

	crypto, err := crypto.NewCrypto(device.Algorithm)
	if err != nil {
		return false, apperrors.WrapError(err, apperrors.InternalError)
	}

	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, apperrors.WrapError(err, apperrors.BadRequest)
	}

	// Verify using the decoded signature
	return crypto.Verify([]byte(dataToBeSigned), decodedSignature, device.PrivateKey)
}
