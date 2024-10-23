package persistence

import (
	"strconv"
	"testing"

	"github.com/chuckiihub/signing-service/domain"
	"github.com/google/uuid"
)

func TestVolatileSignatureSaveAndGet(t *testing.T) {
	memoryStorage := NewVolatileSignatureRepository()
	testSignature := &domain.Signature{
		UUID:       uuid.NewString(),
		SignedData: "signedData",
		Signature:  "s1fn3dD4t4",
		DeviceUUID: "device-uuid",
	}

	testSignature, _ = memoryStorage.Save(testSignature)
	retrievedSignature, err := memoryStorage.FindByUUID(testSignature.UUID)

	if err != nil {
		t.Fatal("Error while recovering Signature")
	}

	if retrievedSignature.SignedData != testSignature.SignedData {
		t.Fatal("Object incorrectly saved")
	}
}

func TestThatSignatureVolatileRepositoryDoesNotReturnPointersToInternalStorageStruct(t *testing.T) {
	const SIGNATURE = "test-Signature"
	const MODIFIED_SIGNATURE = "test-Signature-modified"
	uuid := uuid.NewString()

	memoryStorage := NewVolatileSignatureRepository()

	savedSignature, _ := memoryStorage.Save(&domain.Signature{
		UUID:       uuid,
		SignedData: "signedData",
		Signature:  SIGNATURE,
		DeviceUUID: "device-uuid",
	})

	modifiedSignature, _ := memoryStorage.FindByUUID(savedSignature.UUID)
	modifiedSignature.SignedData = MODIFIED_SIGNATURE

	SignatureInStorage, _ := memoryStorage.FindByUUID(savedSignature.UUID)

	if SignatureInStorage.SignedData == modifiedSignature.SignedData {
		t.Fatal("Modifying objects retrieved by in memory repository modifies objects in internal storage media")
	}
}

func TestListSignaturesWithPagination(t *testing.T) {
	memoryStorage := NewVolatileSignatureRepository()

	// Create and save 3 signatures
	for i := 1; i <= 3; i++ {
		signature := &domain.Signature{
			SignedData: "test-signedData-" + strconv.Itoa(i),
			Signature:  "test-signature-" + strconv.Itoa(i),
			DeviceUUID: "test-device-uuid-" + strconv.Itoa(i),
		}
		_, err := memoryStorage.Save(signature)
		if err != nil {
			t.Fatalf("Error while saving signature %d: %v", i, err)
		}
	}

	// List signatures with page size 2 and check the results
	pageSize := 2
	for page := 1; page < 3; page++ {
		signatures, err := memoryStorage.List(page, pageSize)
		if err != nil {
			t.Fatalf("Error while listing signatures on page %d: %v", page, err)
		}

		expectedSize := pageSize
		if page == 2 {
			expectedSize = 1
		}

		if len(signatures) != expectedSize {
			t.Fatalf("Expected %d signatures on page %d, but got %d", expectedSize, page, len(signatures))
		}

		for i, signature := range signatures {
			expectedDeviceUUID := "test-device-uuid-" + strconv.Itoa((page-1)*pageSize+i+1)
			if signature.DeviceUUID != expectedDeviceUUID {
				t.Fatalf("Expected device UUID %s, but got %s", expectedDeviceUUID, signature.DeviceUUID)
			}
		}
	}
}

func TestListSignaturesPaginationOnEmptyPage(t *testing.T) {
	memoryStorage := NewVolatileSignatureRepository()

	// Create and save 15 devices
	for i := 1; i <= 3; i++ {
		signature := &domain.Signature{
			SignedData: "test-signedData-" + strconv.Itoa(i),
			Signature:  "test-signature-" + strconv.Itoa(i),
			DeviceUUID: "test-device-uuid-" + strconv.Itoa(i),
		}
		_, err := memoryStorage.Save(signature)
		if err != nil {
			t.Fatalf("Error while saving signature %d: %v", i, err)
		}
	}

	signatures, err := memoryStorage.List(20, 10)
	if err != nil {
		t.Fatal("Not expecting error when listing with empty page")
	}

	if len(signatures) != 0 {
		t.Fatal("Expected empty list when listing with invalid page")
	}
}
