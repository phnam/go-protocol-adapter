// Package common provides shared types, constants, and utilities used across the protocol adapter.
package common

import (
	"errors"
	"strings"
)

// Error represents a custom error type for the SDK.
// It contains both an error code and a descriptive message, allowing for
// more structured error handling than standard Go errors.
type Error struct {
	ErrorCode string // Unique identifier for the error type
	Message   string // Human-readable error description
}

// Error implements the error interface by returning a formatted string
// that combines the error code and message with a separator.
func (e Error) Error() string {
	return e.ErrorCode + "//" + e.Message
}

// ToError converts the custom Error type to a standard Go error.
// This is useful when interfacing with code that expects standard errors.
func (e Error) ToError() error {
	return errors.New(e.ErrorCode + "//" + e.Message)
}

// NewError creates a new Error instance with the specified error code and message.
// It returns a pointer to the newly created Error.
func NewError(errorCode string, message string) *Error {
	return &Error{
		ErrorCode: errorCode,
		Message:   message,
	}
}

// ParseError converts a standard Go error into a custom Error type.
// If the error string follows the expected format (code//message), it will
// extract these components. Otherwise, it creates an UNKNOWN_ERROR.
// Returns nil if the input error is nil.
func ParseError(err error) *Error {
	if err == nil {
		return nil
	}
	str := err.Error()
	parts := strings.Split(str, "//")
	if len(parts) != 2 {
		return NewError("UNKNOWN_ERROR", str)
	}
	return NewError(parts[0], parts[1])
}
