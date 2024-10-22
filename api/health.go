package api

import (
	"net/http"

	"github.com/chuckiihub/signing-service/domain"
)

// Health evaluates the health of the service and writes a standardized response.
// The structure that the service answers here goes from this endpoint (API is alive)
// to the SignatureService, DeviceServices and it's DBs.
func (s *Server) Health(response http.ResponseWriter, request *http.Request) {
	signatureService := s.signatureService.CheckHealth()
	deviceService := s.deviceService.CheckHealth()

	health := domain.Health{
		Status:  "pass",
		Version: "v0",
		Services: map[string]domain.ServiceHealth{
			"signature": signatureService,
			"device":    deviceService,
		},
	}

	WriteAPIResponse(response, http.StatusOK, health)
}
