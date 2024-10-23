package persistence

import (
	"strconv"
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

func TestListDevicesWithPagination(t *testing.T) {
	memoryStorage := NewVolatileDeviceRepository()

	// Create and save 15 devices
	for i := 1; i <= 3; i++ {
		device := &domain.Device{
			Label:     "test-device-" + strconv.Itoa(i),
			Algorithm: crypto.SignatureAlgorithmRSA,
		}
		_, err := memoryStorage.Save(device)
		if err != nil {
			t.Fatalf("Error while saving device %d: %v", i, err)
		}
	}

	// List devices with page size 2 and check the results
	pageSize := 2
	for page := 1; page < 3; page++ {
		devices, err := memoryStorage.List(page, pageSize)
		if err != nil {
			t.Fatalf("Error while listing devices on page %d: %v", page, err)
		}

		expectedSize := pageSize
		if page == 2 {
			expectedSize = 1
		}

		if len(devices) != expectedSize {
			t.Fatalf("Expected %d devices on page %d, but got %d", expectedSize, page, len(devices))
		}

		for i, device := range devices {
			expectedLabel := "test-device-" + strconv.Itoa((page-1)*pageSize+i+1)
			if device.Label != expectedLabel {
				t.Fatalf("Expected device label %s, but got %s", expectedLabel, device.Label)
			}
		}
	}
}

func TestListDevicesPaginationOnEmptyPage(t *testing.T) {
	memoryStorage := NewVolatileDeviceRepository()

	// Create and save 15 devices
	for i := 1; i <= 3; i++ {
		device := &domain.Device{
			Label:     "test-device-" + strconv.Itoa(i),
			Algorithm: crypto.SignatureAlgorithmRSA,
		}
		_, err := memoryStorage.Save(device)
		if err != nil {
			t.Fatalf("Error while saving device %d: %v", i, err)
		}
	}

	devices, err := memoryStorage.List(20, 10)
	if err != nil {
		t.Fatal("Not expecting error when listing with empty page")
	}

	if len(devices) != 0 {
		t.Fatal("Expected empty list when listing with invalid page")
	}
}
