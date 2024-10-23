package service

import (
	"log/slog"
	"sync"
)

// If we want to scale horizontally the application, we are gonna be needing
// a external service for locking (like Redis) so 2 different processes don't
// modify fields that should be modified atomically (like signatureCounter or LastSignature).

// Right now, we are only implementing a InMemory Locking Service.

type LockService interface {
	Lock(deviceId string)
	Unlock(deviceId string)
}

type VolatileLockService struct {
	deviceMutexes map[string]*sync.Mutex
	generalMutex  sync.Mutex
}

// To be called while signing data before modiying the `lastSignature`
// and `signatureCounter` fields
func (lockService *VolatileLockService) Lock(deviceId string) {
	lockService.generalMutex.Lock()
	if _, ok := lockService.deviceMutexes[deviceId]; !ok {
		lockService.deviceMutexes[deviceId] = &sync.Mutex{}
	}
	lockService.generalMutex.Unlock()

	lockService.deviceMutexes[deviceId].Lock()
	slog.Debug("lock acquired for device", "deviceId", deviceId)
}

// After successfull save in databaseof the device (and the signature),
// we would need to unlock the device so other threads/processes can modify it.
func (lockService *VolatileLockService) Unlock(deviceId string) {
	mutex, ok := lockService.deviceMutexes[deviceId]
	if !ok {
		slog.Warn("attempting to unlock a non existent lock for device", "deviceId", deviceId)
		return
	}
	mutex.Unlock()
	slog.Debug("lock released for device", "deviceId", deviceId)
}

func NewVolatileLockService() *VolatileLockService {
	return &VolatileLockService{
		deviceMutexes: make(map[string]*sync.Mutex),
	}
}
