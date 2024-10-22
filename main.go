package main

import (
	"log"

	"github.com/chuckiihub/signing-service/api"
	"github.com/chuckiihub/signing-service/persistence"
	"github.com/chuckiihub/signing-service/service"
)

const (
	ListenAddress     = ":8081"
	DevicePageSize    = 20
	SignaturePageSize = 20
)

func main() {
	devicePersistence := persistence.NewVolatileDeviceRepository()
	signaturePersistence := persistence.NewVolatileSignatureRepository()
	lockService := service.NewVolatileLockService()

	deviceService := service.NewDeviceService(devicePersistence, DevicePageSize)
	signatureService := service.NewSignatureService(devicePersistence, signaturePersistence, lockService, SignaturePageSize)

	server := api.NewServer(ListenAddress, deviceService, signatureService)

	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	} else {
		log.Default().Printf("server started on %v", ListenAddress)
	}
}
