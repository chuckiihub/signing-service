package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/chuckiihub/signing-service/api/dto"
	"github.com/gorilla/mux"
)

// The functions here represent the HTTP Transport Layer of the API and are responsible for
// parsing the HTTP request, validating the request, and returning the response.

func (context *Server) SignatureCreate(response http.ResponseWriter, request *http.Request) {
	var creationRequest dto.SignatureCreateRequest
	vars := mux.Vars(request)
	deviceId := vars["deviceId"]

	if deviceId == "" {
		WriteNotFoundError(response)
		return
	}

	err := json.NewDecoder(request.Body).Decode(&creationRequest)
	if err != nil {
		WriteAPIResponse(response, http.StatusBadRequest, err.Error())
		return
	}

	signature, err := context.signatureService.Sign(deviceId, creationRequest.Data)
	if err != nil {
		WriteAppError(response, err)
		return
	}

	signatureResponse := dto.NewSignatureResponseFromSignature(signature)

	WriteAPIResponse(response, http.StatusOK, signatureResponse)
}

func (context *Server) SignatureGet(response http.ResponseWriter, request *http.Request) {
}

// Verifies that a signature is correct given a deviceId, signedData and a Signature given
// by this service
func (context *Server) SignatureVerify(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	deviceId := vars["deviceId"]

	if deviceId == "" {
		WriteNotFoundError(response)
		return
	}

	var verifyRequest dto.SignatureVerifyRequest
	err := json.NewDecoder(request.Body).Decode(&verifyRequest)
	if err != nil {
		WriteAPIResponse(response, http.StatusBadRequest, err.Error())
		return
	}

	if err := context.validator.ValidateSignatureVerificationRequest(verifyRequest); err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, context.validator.GetValidationFailureErrors(err))
		return
	}

	verified, err := context.signatureService.Verify(deviceId, verifyRequest.SignedData, verifyRequest.Signature)
	if err != nil {
		WriteAPIResponse(response, http.StatusTeapot, "invalid")
		return
	}

	if verified {
		WriteAPIResponse(response, http.StatusOK, "valid")
		return
	} else {
		WriteAPIResponse(response, http.StatusTeapot, "invalid")
		return
	}
}

// List services.
func (context *Server) SignatureList(response http.ResponseWriter, request *http.Request) {
	page, err := strconv.Atoi(request.URL.Query().Get("page"))
	if err != nil {
		page = 0
	}

	signatures, err := context.signatureService.List(page)
	if err != nil {
		WriteInternalError(response)
		return
	}

	signaturesResponse := make([]dto.SignatureResponse, 0)
	for _, signature := range signatures {
		signaturesResponse = append(signaturesResponse, *dto.NewSignatureResponseFromSignature(&signature))
	}

	WriteAPIResponse(response, http.StatusOK, signaturesResponse)
}
