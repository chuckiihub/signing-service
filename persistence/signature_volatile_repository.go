package persistence

import (
	"errors"
	"sync"

	"github.com/chuckiihub/signing-service/domain"
)

// In the beginning I've used just a map as a storage structure but
// then I realized that it could be quite not efficient while dealing
// with listing with offsets / pages as I cannot do slices of maps.
// Another way of doing this would be to use a sync.Map.

// The signatureIndexMap is used as an index to the slice of signatures.
// The lock is used when I am adding new signatures so I can update the index.
type SignatureVolatileRepository struct {
	signatureIndexMap map[string]int
	signatures        []domain.Signature
	rwLock            sync.RWMutex
}

func (repository *SignatureVolatileRepository) List(page int, pageSize int) ([]domain.Signature, error) {
	repository.rwLock.RLock()
	defer repository.rwLock.RUnlock()

	if pageSize < 1 {
		return nil, errors.New("batch size cannot be less than 1")
	}

	if page < 0 {
		return nil, errors.New("page cannot be less than 0")
	}

	signatures := make([]domain.Signature, 0, pageSize)

	lowerLimit := (page - 1) * pageSize
	higherLimit := (page) * pageSize

	if lowerLimit > len(repository.signatures) {
		return signatures, nil
	}

	if higherLimit > len(repository.signatures) {
		higherLimit = len(repository.signatures)
	}

	for i := lowerLimit; i < higherLimit; i++ {
		deepCopy := repository.signatures[i]
		signatures = append(signatures, deepCopy)
	}

	return signatures, nil
}

func (repository *SignatureVolatileRepository) Save(signature *domain.Signature) (*domain.Signature, error) {
	repository.rwLock.Lock()
	defer repository.rwLock.Unlock()

	if _, exists := repository.signatureIndexMap[signature.Signature]; exists {
		return signature, errors.New("signature already exists in storage")
	}

	repository.signatures = append(repository.signatures, *signature)
	repository.signatureIndexMap[signature.UUID] = len(repository.signatures) - 1

	return signature, nil
}

func (repository *SignatureVolatileRepository) FindByUUID(uuid string) (*domain.Signature, error) {
	repository.rwLock.RLock()
	defer repository.rwLock.RUnlock()

	if index, exists := repository.signatureIndexMap[uuid]; exists {
		deepCopy := repository.signatures[index]
		return &deepCopy, nil
	}

	return nil, nil
}

func (repository *SignatureVolatileRepository) CheckHealth() domain.PersistenceHealth {
	return domain.PersistenceHealth{Status: domain.HealthStatusPass}
}
