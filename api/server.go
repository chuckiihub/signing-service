package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	apperrors "github.com/chuckiihub/signing-service/errors"
	"github.com/chuckiihub/signing-service/service"
	"github.com/gorilla/mux"
)

// Response is the generic API response container.
type Response struct {
	Data interface{} `json:"data"`
}

// ErrorResponse is the generic error API response container.
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// Server manages HTTP requests and dispatches them to the appropriate services.
type Server struct {
	listenAddress    string
	deviceService    service.DeviceService
	signatureService service.SignatureService
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string, deviceService service.DeviceService, signatureService service.SignatureService) *Server {
	return &Server{
		listenAddress:    listenAddress,
		deviceService:    deviceService,
		signatureService: signatureService,
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	router := mux.NewRouter()

	router.HandleFunc("/api/v0/docs", s.ServeDocs).Methods("GET")

	router.HandleFunc("/api/v0/health", s.Health)

	router.HandleFunc("/api/v0/device", s.DeviceCreate).Methods("POST")
	router.HandleFunc("/api/v0/device/{uuid}", s.DeviceGet).Methods("GET")
	router.HandleFunc("/api/v0/device", s.DeviceList).Methods("GET")

	router.HandleFunc("/api/v0/device/{deviceId}/sign", s.SignatureCreate).Methods("POST")
	// Using post as the signedData might be large
	router.HandleFunc("/api/v0/device/{deviceId}/verify", s.SignatureVerify).Methods("POST")

	router.HandleFunc("/api/v0/signature/{signature}", s.SignatureGet).Methods("GET")
	router.HandleFunc("/api/v0/signature", s.SignatureList).Methods("GET")

	return http.ListenAndServe(s.listenAddress, router)
}

func (s *Server) ServeDocs(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/docs/api.html")
}

// Handy function, for example, when JSON body is not valid
func WriteInvalidRequestBodyError(w http.ResponseWriter) {
	WriteErrorResponse(w, http.StatusBadRequest, []string{http.StatusText(http.StatusBadRequest)})
}

// WriteInternalError writes a default internal error message as an HTTP response.
func WriteInternalError(w http.ResponseWriter) {
	WriteAPIResponse(w, http.StatusInternalServerError, []byte(http.StatusText(http.StatusInternalServerError)))
}

// WriteErrorResponse takes an HTTP status code and a slice of errors
// and writes those as an HTTP error response in a structured format.
func WriteErrorResponse(w http.ResponseWriter, code int, errors []string) {
	w.WriteHeader(code)

	errorResponse := ErrorResponse{
		Errors: errors,
	}

	bytes, err := json.Marshal(errorResponse)
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}

// WriteAPIResponse takes an HTTP status code and a generic data struct
// and writes those as an HTTP response in a structured format.
func WriteAPIResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)

	response := Response{
		Data: data,
	}

	bytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}

// WriteAPIResponse takes an HTTP status code and a generic data struct
// and writes those as an HTTP response in a structured format.
func WriteNotFoundError(w http.ResponseWriter) {
	WriteErrorResponse(w, http.StatusNotFound, []string{http.StatusText(http.StatusNotFound)})
}

// The Device and Signature services might return different types of errors
// and as those are not responsible for the transport layer, I am wrapping
// the errors so the caller (HandlerFunc) knows what type of Status Code should
// return.
func WriteAppError(w http.ResponseWriter, err error) {
	var appErr *apperrors.AppError

	statusCode := http.StatusInternalServerError
	if errors.As(err, &appErr) {
		statusCode = appErr.Type
	}

	content := []string{http.StatusText(http.StatusBadRequest)}
	if err.Error() != "" {
		content = []string{err.Error()}
	}

	if statusCode == http.StatusInternalServerError {
		slog.Error("Unhandled internal error", "error", err.Error())
	}

	WriteErrorResponse(w, statusCode, content)
}
