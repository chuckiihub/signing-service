package main

import (
	"log/slog"
	"os"

	"github.com/chuckiihub/signing-service/api"
	"github.com/chuckiihub/signing-service/config"
	"github.com/chuckiihub/signing-service/persistence"
	"github.com/chuckiihub/signing-service/service"
)

func configureLogging() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)
}

func main() {
	configureLogging()

	devicePersistence := persistence.NewVolatileDeviceRepository()
	signaturePersistence := persistence.NewVolatileSignatureRepository()
	lockService := service.NewVolatileLockService()

	deviceService := service.NewDeviceService(devicePersistence, config.ListPageSize)
	signatureService := service.NewSignatureService(devicePersistence, signaturePersistence, lockService, config.ListPageSize)

	listenAddress := config.GetListenAddress(config.DefaultListenAddress)
	server := api.NewServer(listenAddress, deviceService, signatureService)

	if err := server.Run(); err != nil {
		slog.Error("could not start server", "port", listenAddress, "error", err.Error())
	} else {
		slog.Info("server started", "port", listenAddress)
	}
}
