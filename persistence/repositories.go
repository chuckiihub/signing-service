package persistence

import (
	"github.com/chuckiihub/signing-service/domain"
)

type DevicePersistance interface {
	Save(device *domain.Device) (*domain.Device, error)
	FindByUUID(UUID string) (*domain.Device, error)
	List(page int, pageSize int) ([]domain.Device, error)
	CheckHealth() domain.PersistenceHealth
}

func NewVolatileDeviceRepository() *VolatileDeviceRepository {
	return &VolatileDeviceRepository{
		devices:   make([]domain.Device, 0),
		uuidIndex: make(map[string]int),
	}
}

type SignaturePersistance interface {
	List(offset int, pageSize int) ([]domain.Signature, error)
	Save(signature *domain.Signature) (*domain.Signature, error)
	FindBySignature(signature string) (*domain.Signature, error)
	CheckHealth() domain.PersistenceHealth
}

func NewVolatileSignatureRepository() *SignatureVolatileRepository {
	return &SignatureVolatileRepository{
		signatureIndexMap: make(map[string]int, 0),
		signatures:        make([]domain.Signature, 0),
	}
}

type PersistenceHealthCheck interface {
	CheckHealth() domain.PersistenceHealth
}
