package service

import (
	"encoding/base64"
	"testing"

	"github.com/chuckiihub/signing-service/crypto"
	"github.com/chuckiihub/signing-service/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDevicePersistence struct {
	mock.Mock
}

func (m *MockDevicePersistence) Save(device *domain.Device) (*domain.Device, error) {
	args := m.Called(device)
	return args.Get(0).(*domain.Device), args.Error(1)
}

func (m *MockDevicePersistence) FindByUUID(uuid string) (*domain.Device, error) {
	args := m.Called(uuid)
	return args.Get(0).(*domain.Device), args.Error(1)
}

func (m *MockDevicePersistence) List(page int, pageSize int) ([]domain.Device, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]domain.Device), args.Error(1)
}

func (m *MockDevicePersistence) CheckHealth() domain.PersistenceHealth {
	args := m.Called()
	return args.Get(0).(domain.PersistenceHealth)
}

func TestCreate(t *testing.T) {
	mockPersistence := new(MockDevicePersistence)
	deviceService := &DeviceServiceImplementation{
		persistence: mockPersistence,
		pageSize:    10,
	}

	uuid := "test-uuid"
	algorithm := crypto.SignatureAlgorithmRSA
	label := "someLabel"
	privateKey := []byte("private")
	publicKey := []byte("public")

	expectedDevice := &domain.Device{
		UUID:          uuid,
		Algorithm:     algorithm,
		Label:         label,
		LastSignature: base64.StdEncoding.EncodeToString([]byte(uuid)),
		PrivateKey:    privateKey,
		PublicKey:     publicKey,
	}

	mockPersistence.On("Save", mock.AnythingOfType("*domain.Device")).Return(expectedDevice, nil)

	device, err := deviceService.Create(uuid, algorithm, label)

	assert.NoError(t, err)
	assert.NotNil(t, device)
	assert.Equal(t, uuid, device.UUID)
	assert.Equal(t, algorithm, device.Algorithm)
	assert.Equal(t, label, device.Label)
	assert.Equal(t, privateKey, device.PrivateKey)
	assert.Equal(t, publicKey, device.PublicKey)

	mockPersistence.AssertExpectations(t)
}

func TestFind(t *testing.T) {
	mockPersistence := new(MockDevicePersistence)
	deviceService := &DeviceServiceImplementation{
		persistence: mockPersistence,
		pageSize:    10,
	}

	uuid := "test-uuid"
	expectedDevice := &domain.Device{UUID: uuid}

	mockPersistence.On("FindByUUID", uuid).Return(expectedDevice, nil)

	device, err := deviceService.Get(uuid)

	assert.NoError(t, err)
	assert.NotNil(t, device)
	assert.Equal(t, uuid, device.UUID)

	mockPersistence.AssertExpectations(t)
}

func TestList(t *testing.T) {
	pageSize := 5
	page := 1

	mockPersistence := new(MockDevicePersistence)
	deviceService := &DeviceServiceImplementation{
		persistence: mockPersistence,
		pageSize:    pageSize,
	}

	expectedDevices := []domain.Device{{UUID: "1"}, {UUID: "2"}}

	mockPersistence.On("List", page, pageSize).Return(expectedDevices, nil)

	devices, err := deviceService.List(page)

	assert.NoError(t, err)
	assert.NotNil(t, devices)
	assert.Len(t, devices, 2)

	mockPersistence.AssertExpectations(t)
}
