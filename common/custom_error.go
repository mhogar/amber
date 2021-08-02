package common

import "errors"

const (
	ErrorTypeNone     = iota
	ErrorTypeInternal = iota
	ErrorTypeClient   = iota
)

// CustomError is an error with an added type field to determine how it should be handled
type CustomError struct {
	error
	Type int
}

// NoError returns a CustomError with type ErrorTypeNone
func NoError() CustomError {
	return CustomError{
		error: nil,
		Type:  ErrorTypeNone,
	}
}

// InternalError returns a CustomError with type ErrorTypeInternal and an internal error message
func InternalError() CustomError {
	return CustomError{
		error: errors.New("an internal error occurred"),
		Type:  ErrorTypeInternal,
	}
}

// ClientError returns a CustomError with type ErrorTypeClient and the provided message
func ClientError(message string) CustomError {
	return CustomError{
		error: errors.New(message),
		Type:  ErrorTypeClient,
	}
}

// OAuthCustomError is a CustomError with an added error name field to determine the oauth error name
type OAuthCustomError struct {
	CustomError
	ErrorName string
}

// OAuthNoError returns an OAuthCustomError with type ErrorTypeNone
func OAuthNoError() OAuthCustomError {
	return OAuthCustomError{
		CustomError: NoError(),
		ErrorName:   "",
	}
}

// OAuthInternalError returns an OAuthCustomError with type ErrorTypeInternal and an internal error message
func OAuthInternalError() OAuthCustomError {
	return OAuthCustomError{
		CustomError: InternalError(),
		ErrorName:   "",
	}
}

// ClientError returns a CustomError with type ErrorTypeClient and the provided error name and message
func OAuthClientError(errorName string, message string) OAuthCustomError {
	return OAuthCustomError{
		CustomError: ClientError(message),
		ErrorName:   errorName,
	}
}
