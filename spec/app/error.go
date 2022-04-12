package app

import "errors"

// NotFound means that the resource has not been found.
func NotFound(err error) bool {
	var notAllowed interface {
		NotFound() bool
	}

	return errors.As(err, &notAllowed) && notAllowed.NotFound()
}

// Forbidden means that the user is authenticated but the authorization has failed. This is usually permanent,
// because the access control mechanism forbids the access to that resource.
// This is equal to http.StatusForbidden (403).
func Forbidden(err error) bool {
	var notAllowed interface {
		Forbidden() bool
	}

	return errors.As(err, &notAllowed) && notAllowed.Forbidden()
}

// Unauthenticated is equal to http.StatusUnauthorized (401) and means that no authentication is available but
// is generally required. This is likely solved by performing a login and passing a valid token.
func Unauthenticated(err error) bool {
	var notAllowed interface {
		Unauthenticated() bool
	}

	return errors.As(err, &notAllowed) && notAllowed.Unauthenticated()
}

func InternalServerError(err error) bool {
	var notAllowed interface {
		InternalServerError() bool
	}

	return errors.As(err, &notAllowed) && notAllowed.InternalServerError()
}

func ProtocolError(err error) {

}

// ValidationError describes a condition where a validation has failed.
type ValidationError struct {
	Message string
	Cause   error
}

func (e ValidationError) Unwrap() error {
	return e.Cause
}

func (e ValidationError) Error() string {
	return e.Message
}

func (e ValidationError) FailedValidation() bool {
	return true
}

func ValidationFailed(err error) (bool, string) {
	var e interface {
		FailedValidation() bool
	}

	if errors.As(err, &e) && e.FailedValidation() {
		return true, err.Error()
	}

	return false, ""
}
