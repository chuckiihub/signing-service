package persistence

import (
	"testing"

	"github.com/chuckiihub/signing-service/crypto"
	"github.com/chuckiihub/signing-service/domain"
)

func TestSaveWithoutUUID(t *testing.T) {
	memoryStorage := NewVolatileDeviceRepository()
	device := domain.Device{
		Label:     "test-device",
		Algorithm: crypto.SignatureAlgorithmRSA,
	}

	savedDevice, err := memoryStorage.Save(&device)

	if err != nil {
		t.Fatalf(`Error while saving %v`, err)
	}

	if savedDevice.UUID == "" {
		t.Fatalf(`Saved device has no UUID`)
	}

	if device.UUID == "" {
		t.Fatalf(`Original device has no UUID`)
	}
}

func TestSaveAndGet(t *testing.T) {
	memoryStorage := NewVolatileDeviceRepository()
	device := &domain.Device{
		Label:     "test-device",
		Algorithm: crypto.SignatureAlgorithmRSA,
	}

	savedDevice, _ := memoryStorage.Save(device)
	device, err := memoryStorage.FindByUUID(savedDevice.UUID)

	if err != nil {
		t.Fatal("Error while recovering device")
	}

	if device.Label != savedDevice.Label {
		t.Fatal("Object incorrectly saved")
	}
}

func TestThatRepositoryDoesNotReturnPointersToInternalStorageStruct(t *testing.T) {
	const LABEL = "test-device"
	const MODIFIED_LABEL = "test-device-modified"

	memoryStorage := NewVolatileDeviceRepository()

	savedDevice, _ := memoryStorage.Save(&domain.Device{
		Label:     LABEL,
		Algorithm: crypto.SignatureAlgorithmRSA,
	})

	modifiedDevice, _ := memoryStorage.FindByUUID(savedDevice.UUID)
	modifiedDevice.Label = MODIFIED_LABEL

	deviceInStorage, _ := memoryStorage.FindByUUID(savedDevice.UUID)

	if deviceInStorage.Label == modifiedDevice.Label {
		t.Fatal("Modifying objects retrieved by in memory repository modifies objects in internal storage media")
	}
}
