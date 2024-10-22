package apperrors

// Did not use the standard http error codes as I did not wanted to
// import the package unnecesarily.
const (
	InternalError = 500
	NotFound      = 404
	BadRequest    = 400
)

type AppError struct {
	Type int

	Err error
}

func (e AppError) Error() string {
	return e.Err.Error()
}

func WrapError(err error, code int) AppError {
	return AppError{
		Type: code,
		Err:  err,
	}
}

func (e *AppError) Unwrap() error {
	return e.Err
}
