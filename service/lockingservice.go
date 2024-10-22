package service

import "sync"

// If we want to scale horizontally the application, we are gonna be needing
// a external service for locking (like Redis) so 2 different processes don't
// modify fields that should be modified atomically (like signatureCounter or LastSignature).

// Right now, we are only implementing a InMemory Locking Service.

type LockService interface {
	Lock(deviceId string)
	Unlock(deviceId string)
}

type VolatileLockService struct {
	deviceMutexes map[string]*sync.RWMutex
	generalMutex  sync.RWMutex
}

// To be called while signing data before modiying the `lastSignature`
// and `signatureCounter` fields
func (lockService *VolatileLockService) Lock(deviceId string) {
	lockService.generalMutex.Lock()
	defer lockService.generalMutex.Unlock()

	if _, ok := lockService.deviceMutexes[deviceId]; !ok {
		lockService.deviceMutexes[deviceId] = &sync.RWMutex{}
	}
	lockService.deviceMutexes[deviceId].Lock()
}

// After successfull save in databaseof the device (and the signature),
// we would need to unlock the device so other threads/processes can modify it.
func (lockService *VolatileLockService) Unlock(deviceId string) {
	lockService.generalMutex.Lock()
	defer lockService.generalMutex.Unlock()

	lockService.deviceMutexes[deviceId].Unlock()
}

func NewVolatileLockService() *VolatileLockService {
	return &VolatileLockService{
		deviceMutexes: make(map[string]*sync.RWMutex),
	}
}
