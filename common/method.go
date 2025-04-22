// Package common provides shared types, constants, and utilities used across the protocol adapter.
package common

// MethodValue represents a single HTTP method value.
// It encapsulates the string representation of an HTTP method.
type MethodValue struct {
	Value string // The string representation of the HTTP method
}

// MethodEnum defines a structure containing all supported HTTP methods.
// Each field is a pointer to a MethodValue representing a specific HTTP method.
type MethodEnum struct {
	GET     *MethodValue // HTTP GET method
	QUERY   *MethodValue // Custom QUERY method
	POST    *MethodValue // HTTP POST method
	PUT     *MethodValue // HTTP PUT method
	PATCH   *MethodValue // HTTP PATCH method
	DELETE  *MethodValue // HTTP DELETE method
	OPTIONS *MethodValue // HTTP OPTIONS method
}

// APIMethod is a published enum containing predefined HTTP method values.
// It provides a convenient way to reference standard HTTP methods throughout the application.
var APIMethod = MethodEnum{
	GET:     &MethodValue{Value: "GET"},
	QUERY:   &MethodValue{Value: "QUERY"},
	POST:    &MethodValue{Value: "POST"},
	PUT:     &MethodValue{Value: "PUT"},
	PATCH:   &MethodValue{Value: "PATCH"},
	DELETE:  &MethodValue{Value: "DELETE"},
	OPTIONS: &MethodValue{Value: "OPTIONS"},
}
