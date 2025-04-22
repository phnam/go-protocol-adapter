// Package common provides shared types, constants, and utilities used across the protocol adapter.
package common

import (
	"errors"
	"reflect"
	"strings"
)

// APIResponse represents a standardized response object with JSON format.
// It provides a consistent structure for all API responses, including success and error cases.
// The generic type parameter T allows for type-safe data handling.
type APIResponse[T any] struct {
	Status    string            `json:"status"`               // Response status (e.g., "OK", "ERROR")
	Data      []T               `json:"data,omitempty"`       // Array of response data objects
	Message   string            `json:"message"`              // Human-readable message
	ErrorCode string            `json:"error_code,omitempty"` // Error code in case of failure
	Total     int64             `json:"total,omitempty"`      // Total count of items (for pagination)
	Headers   map[string]string `json:"headers,omitempty"`    // Response headers
}

// ToAnyResponse converts a typed APIResponse to a generic APIResponse with 'any' type.
// This is useful when you need to handle responses of different types uniformly.
func (resp *APIResponse[T]) ToAnyResponse() *APIResponse[any] {
	arr := []any{}
	for _, v := range resp.Data {
		arr = append(arr, v)
	}
	return &APIResponse[any]{
		Status:    resp.Status,
		Data:      arr,
		Message:   resp.Message,
		ErrorCode: resp.ErrorCode,
		Total:     resp.Total,
		Headers:   resp.Headers,
	}
}

// FromError converts a standard error or custom Error into an APIResponse.
// It analyzes the error type and content to determine the appropriate response status and error code.
// If the error is nil, it returns a success response.
func FromError(err error) *APIResponse[any] {

	var e Error
	if errors.As(err, &e) {
		// Handle custom Error type
		if e.ErrorCode == "NOT_FOUND" {
			return NewErrorResponse(APIStatus.NotFound, e.ErrorCode, e.Message)
		}
		return NewErrorResponse(APIStatus.NotFound, e.ErrorCode, e.Message)
	}

	if err != nil {
		// Parse error string in format "CODE//MESSAGE"
		msgParts := strings.Split(err.Error(), "//")
		if len(msgParts) != 2 {
			// Handle non-standard error format
			return NewErrorResponse(APIStatus.Error, "INTERNAL_SERVER_ERROR", err.Error())
		}
		errorCode := msgParts[0]

		// Map error codes to appropriate response statuses
		if errorCode == "NOT_FOUND" {
			return NewErrorResponse(APIStatus.NotFound, errorCode, msgParts[1])
		}
		if strings.HasPrefix(errorCode, "INVALID") {
			return NewErrorResponse(APIStatus.Invalid, errorCode, msgParts[1])
		}
		if strings.HasPrefix(errorCode, "EXISTED") {
			return NewErrorResponse(APIStatus.Existed, errorCode, msgParts[1])
		}
		if strings.HasPrefix(errorCode, "FORBIDDEN") {
			return NewErrorResponse(APIStatus.Forbidden, errorCode, msgParts[1])
		}
		if strings.HasPrefix(errorCode, "UNAUTHORIZED") {
			return NewErrorResponse(APIStatus.Unauthorized, errorCode, msgParts[1])
		}
		if strings.HasPrefix(errorCode, "REDIRECTED") {
			return NewErrorResponse(APIStatus.Redirected, errorCode, msgParts[1])
		}

		// Default error response
		return NewErrorResponse(APIStatus.Error, errorCode, msgParts[1])
	}
	// No error, return success response
	return NewOkResponse(nil, "Success")
}

// NewAPIResponse creates a new APIResponse with the specified parameters.
// It handles both array and single-item data by ensuring the Data field is always an array.
// If data is already a slice, it's used directly; otherwise, it's wrapped in a single-element array.
func NewAPIResponse(status string, data []any, message string, errorCode string, total int64, headers map[string]string) *APIResponse[any] {
	// Check if data is a slice with Reflect
	if data == nil || reflect.TypeOf(data).Kind() == reflect.Slice {
		return &APIResponse[any]{
			Status:    status,
			Data:      data,
			Message:   message,
			ErrorCode: errorCode,
			Total:     total,
			Headers:   headers,
		}
	}

	// Wrap single item in an array
	return &APIResponse[any]{
		Status:    status,
		Data:      []interface{}{data},
		Message:   message,
		ErrorCode: errorCode,
		Total:     total,
		Headers:   headers,
	}
}

// NewErrorResponse creates an error response with the specified status, error code, and message.
// It returns an APIResponse with no data, focusing on the error information.
func NewErrorResponse(status string, errorCode string, message string) *APIResponse[any] {
	return &APIResponse[any]{
		Status:    status,
		Message:   message,
		ErrorCode: errorCode,
	}
}

// NewOkResponse creates a success response with the specified data and message.
// It automatically sets the status to APIStatus.Ok and handles both array and single-item data.
func NewOkResponse(data []any, message string) *APIResponse[any] {

	// Check if data is a slice with Reflect
	if data == nil || reflect.TypeOf(data).Kind() == reflect.Slice {
		return &APIResponse[any]{
			Status:  APIStatus.Ok,
			Data:    data,
			Message: message,
		}
	}

	// Wrap single item in an array
	return &APIResponse[any]{
		Status:  APIStatus.Ok,
		Data:    []interface{}{data},
		Message: message,
	}
}

// StatusEnum defines a structure containing all possible API response status values.
// These statuses are used to indicate the result of an API operation.
type StatusEnum struct {
	Ok           string // Successful operation
	Error        string // General error
	Invalid      string // Invalid input or request
	NotFound     string // Requested resource not found
	Forbidden    string // Access denied
	Existed      string // Resource already exists
	Unauthorized string // Authentication required
	Redirected   string // Request redirected
}

// APIStatus is a published enum containing predefined status values.
// It provides a consistent way to set response statuses throughout the application.
var APIStatus = &StatusEnum{
	Ok:           "OK",
	Error:        "ERROR",
	Invalid:      "INVALID",
	NotFound:     "NOT_FOUND",
	Forbidden:    "FORBIDDEN",
	Existed:      "EXISTED",
	Unauthorized: "UNAUTHORIZED",
	Redirected:   "REDIRECTED",
}
