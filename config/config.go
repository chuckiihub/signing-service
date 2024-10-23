package config

import (
	"log/slog"
	"os"
)

// I've moved this here as when the service grows it might handle configuration
// in a more complex way. Could be fetched from a yaml file as well or environment
// variables so, for example, configuration in pods/containers is easier.

// For now we keep it as it was.
const (
	DefaultListenAddress = ":8081"
	DefaultLogLevel      = slog.LevelInfo
	ListPageSize         = 2
)

// tries to fetch listen address from environment variable, if not found, returns default
func GetListenAddress(defaultAddress string) string {
	listenAddress := os.Getenv("SIGNING_SERVICE_LISTEN_ADDRESS")
	if listenAddress == "" {
		return defaultAddress
	}
	return listenAddress
}

// tries to fetch log level from environment variable, if not found, returns default
func GetLogLevel() slog.Level {
	logLevel := os.Getenv("SIGNING_SERVICE_LOG_LEVEL")

	switch logLevel {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return DefaultLogLevel
	}
}
