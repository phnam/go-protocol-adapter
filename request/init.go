// Package request provides interfaces and implementations for handling API requests across different protocols.
// It abstracts the underlying transport mechanisms (HTTP, Thrift, etc.) to provide a unified API.
package request

import (
	"github.com/phnam/go-protocol-adapter/common"
)

// APIRequest defines the interface for all request types in the application.
// It provides protocol-agnostic methods to access request data regardless of the underlying transport.
type APIRequest interface {
	// GetPath returns the request path/endpoint
	GetPath() string

	// GetMethod returns the HTTP method or equivalent operation type
	GetMethod() *common.MethodValue

	// GetParam retrieves a single query parameter by name
	GetParam(string) string

	// GetParams returns all query parameters as a map
	GetParams() map[string]string

	// GetHeader retrieves a single header value by name
	GetHeader(string) string

	// GetHeaders returns all headers as a map
	GetHeaders() map[string]string

	// ParseBody unmarshals the request body into the provided interface
	ParseBody(interface{}) error

	// GetContentText returns the raw request body as a string
	GetContentText() string

	// GetAttribute retrieves a context attribute by name
	GetAttribute(string) interface{}

	// SetAttribute stores a context attribute
	SetAttribute(string, interface{})

	// SetVar sets a path variable/parameter
	SetVar(string, string)

	// GetVar retrieves a path variable/parameter by name
	GetVar(string) string

	// GetIP returns the client's IP address
	GetIP() string
}
