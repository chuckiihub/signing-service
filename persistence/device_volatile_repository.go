package persistence

import (
	"errors"
	"sync"

	"github.com/chuckiihub/signing-service/domain"
	"github.com/google/uuid"
)

// In the beginning I've used just a map as a storage structure but
// then I realized that it could be quite not efficient while dealing
// with listing with offsets / pages as I cannot do slices of maps.
// Another way of doing this would be to use a sync.Map.

// The uuidIndex is used as an index to the slice of devices.
// The lock is used when I am adding new devices so I can update the index.
type VolatileDeviceRepository struct {
	uuidIndex map[string]int
	devices   []domain.Device
	rwMutex   sync.RWMutex
}

func (repository *VolatileDeviceRepository) Save(device *domain.Device) (*domain.Device, error) {
	repository.rwMutex.Lock()
	defer repository.rwMutex.Unlock()

	if device.UUID == "" {
		device.UUID = uuid.NewString()
	}

	if index, exists := repository.uuidIndex[device.UUID]; exists {
		repository.devices[index] = *device
	} else {
		repository.devices = append(repository.devices, *device)
		repository.uuidIndex[device.UUID] = len(repository.devices) - 1
	}

	return device, nil
}

func (repository *VolatileDeviceRepository) FindByUUID(UUID string) (*domain.Device, error) {
	repository.rwMutex.RLock()
	defer repository.rwMutex.RUnlock()

	if deviceIndex, exists := repository.uuidIndex[UUID]; exists {
		deepCopy := repository.devices[deviceIndex]
		return &deepCopy, nil
	}

	return nil, nil
}

func (repository *VolatileDeviceRepository) List(page int, batchSize int) ([]domain.Device, error) {
	repository.rwMutex.RLock()
	defer repository.rwMutex.RUnlock()
	devices := make([]domain.Device, 0, batchSize)

	if batchSize < 1 {
		return nil, errors.New("batch size cannot be less than 1")
	}

	if page < 0 {
		return nil, errors.New("page cannot be less than 0")
	}
	lowerLimit := page * batchSize
	higherLimit := (page + 1) * batchSize

	if lowerLimit > len(repository.devices) {
		return devices, nil
	}

	if higherLimit > len(repository.devices) {
		higherLimit = len(repository.devices)
	}

	for i := lowerLimit; i < higherLimit; i++ {
		deepCopy := repository.devices[i]
		devices = append(devices, deepCopy)
	}

	return devices, nil
}

func (repository *VolatileDeviceRepository) CheckHealth() domain.PersistenceHealth {
	return domain.PersistenceHealth{Status: domain.HealthStatusPass}
}
