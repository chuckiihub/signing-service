package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/chuckiihub/signing-service/api/dto"
	"github.com/chuckiihub/signing-service/api/validation"
	"github.com/gorilla/mux"
)

// The functions here represent the HTTP Transport Layer of the API and are responsible for
// parsing the HTTP request, validating the request, and returning the response.

func (context *Server) DeviceCreate(response http.ResponseWriter, request *http.Request) {
	var creationRequest dto.DeviceCreationRequest

	err := json.NewDecoder(request.Body).Decode(&creationRequest)
	if err != nil {
		WriteInvalidRequestBodyError(response)
		return
	}

	validator := validation.NewRequestValidator()
	if err := validator.Validate(creationRequest); err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, validator.GetValidationFailureErrors(err))
		return
	}

	signatureAlgorithm, err := creationRequest.GetSignatureAlgorithm()
	if err != nil {
		// this is actually handled by the validator.
		WriteErrorResponse(response, http.StatusBadRequest, []string{"not supported algorithm"})
		return
	}

	device, err := context.deviceService.Create(
		signatureAlgorithm,
		creationRequest.Label,
	)

	if err != nil {
		WriteAppError(response, err)
		return
	}

	WriteAPIResponse(response, http.StatusCreated, dto.NewDeviceResponse(device))
}

func (context *Server) DeviceGet(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	uuid := vars["uuid"]

	if uuid == "" {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"UUID is required",
		})
		return
	}

	device, err := context.deviceService.Get(uuid)
	if err != nil {
		WriteAppError(response, err)
		return
	}

	if device == nil {
		WriteNotFoundError(response)
		return
	}

	WriteAPIResponse(response, http.StatusOK, dto.NewDeviceResponse(device))
}

func (context *Server) DeviceList(response http.ResponseWriter, request *http.Request) {
	pageString := request.URL.Query().Get("page")
	page, err := strconv.Atoi(pageString)

	if err != nil || page < 1 || pageString == "" {
		page = 1
	}

	devices, err := context.deviceService.List(page)
	if err != nil {
		WriteAppError(response, err)
		return
	}

	devicesResponse := make([]dto.DeviceResponse, 0)
	for _, device := range devices {
		devicesResponse = append(devicesResponse, dto.NewDeviceResponse(&device))
	}

	WriteAPIResponse(response, http.StatusOK, devicesResponse)
}
