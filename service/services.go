package service

import (
	"github.com/chuckiihub/signing-service/crypto"
	"github.com/chuckiihub/signing-service/domain"
	"github.com/chuckiihub/signing-service/persistence"
)

type SignatureService interface {
	Sign(deviceId string, dataToBeSigned string) (*domain.Signature, error)
	Verify(deviceId string, dataToBeSigned string, signature string) (bool, error)
	Get(signature string) (*domain.Signature, error)
	List(page int) ([]domain.Signature, error)
	CheckHealth() domain.ServiceHealth
}

type DeviceService interface {
	Create(uuid string, algorithm crypto.SignatureAlgorithm, label string) (*domain.Device, error)
	Get(uuid string) (*domain.Device, error)
	List(page int) ([]domain.Device, error)
	CheckHealth() domain.ServiceHealth
}

func NewSignatureService(
	dDB persistence.DevicePersistance,
	sDB persistence.SignaturePersistance,
	l LockService,
	pageSize int) SignatureService {
	return &SignatureServiceImplementation{
		devicePersistence:    dDB,
		signaturePersistence: sDB,
		lockService:          l,
		pageSize:             pageSize,
	}
}

func NewDeviceService(
	persistence persistence.DevicePersistance,
	pageSize int,
) DeviceService {
	return &DeviceServiceImplementation{
		persistence: persistence,
		pageSize:    pageSize,
	}
}
