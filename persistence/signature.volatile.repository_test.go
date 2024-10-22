package persistence

import (
	"testing"

	"github.com/chuckiihub/signing-service/domain"
)

func TestVolatileSignatureSaveAndGet(t *testing.T) {
	memoryStorage := NewVolatileSignatureRepository()
	testSignature := &domain.Signature{
		SignedData: "signedData",
		Signature:  "s1fn3dD4t4",
		DeviceUUID: "device-uuid",
	}

	testSignature, _ = memoryStorage.Save(testSignature)
	retrievedSignature, err := memoryStorage.FindBySignature(testSignature.Signature)

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

	memoryStorage := NewVolatileSignatureRepository()

	savedSignature, _ := memoryStorage.Save(&domain.Signature{
		SignedData: "signedData",
		Signature:  SIGNATURE,
		DeviceUUID: "device-uuid",
	})

	modifiedSignature, _ := memoryStorage.FindBySignature(savedSignature.Signature)
	modifiedSignature.SignedData = MODIFIED_SIGNATURE

	SignatureInStorage, _ := memoryStorage.FindBySignature(savedSignature.Signature)

	if SignatureInStorage.SignedData == modifiedSignature.SignedData {
		t.Fatal("Modifying objects retrieved by in memory repository modifies objects in internal storage media")
	}
}
