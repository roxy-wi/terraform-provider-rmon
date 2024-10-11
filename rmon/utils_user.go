package rmon

import (
	"fmt"
	"net/http"
	"net/mail"
)

// Utility function to check if the error is a 404 not found error
func isNotFound(err error) bool {
	if httpErr, ok := err.(*httpError); ok {
		return httpErr.StatusCode == http.StatusNotFound
	}
	return false
}

// Define the HTTPError struct and methods
type httpError struct {
	StatusCode int
	Err        error
}

func (e *httpError) Error() string {
	return e.Err.Error()
}

// Utility function to validate email format
func validateEmail(val interface{}, key string) (warns []string, errs []error) {
	_, err := mail.ParseAddress(val.(string))
	if err != nil {
		errs = append(errs, fmt.Errorf("%q must be a valid email address: %v", key, err))
	}
	return
}
