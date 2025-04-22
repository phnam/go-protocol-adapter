// Package server provides implementations for different protocol servers (HTTP and Thrift).
// It defines a common Server interface and protocol-specific implementations.
package server

import (
	"net/http"
	"sync"

	"github.com/phnam/go-protocol-adapter/common"
)

// idCounter is used to generate unique IDs for server instances
var idCounter = 0

// ServerConfig defines configuration options for server instances.
// It contains settings that apply to all server types as well as
// protocol-specific settings.
type ServerConfig struct {
	// Protocol specifies which protocol implementation to use ("HTTP" or "THRIFT")
	Protocol string

	// HideFuncName determines whether function names should be included in response headers
	HideFuncName bool

	// BufferSize specifies the buffer size in bytes for Thrift server transport
	BufferSize int

	// MessageSize specifies the maximum message size in bytes for Thrift server
	MessageSize int32
}

// Server defines the common interface for all protocol server implementations.
// It provides methods for configuring routes, handling requests, and starting the server.
type Server interface {
	// PreRequest registers a handler function that will be executed before every request.
	// This can be used for authentication, logging, or other cross-cutting concerns.
	PreRequest(Handler) error

	// SetHandler registers a handler function for a specific HTTP method and path.
	// The method parameter specifies the HTTP method (GET, POST, etc.)
	// The path parameter specifies the URL path to match
	// The fn parameter is the handler function to execute when the route is matched
	SetHandler(*common.MethodValue, string, Handler) error

	// Expose sets the port number that the server will listen on
	Expose(int)

	// Start begins listening for incoming requests on the configured port.
	// It takes a WaitGroup parameter to allow the caller to wait for the server to exit.
	Start(*sync.WaitGroup)

	// GetHostname returns the hostname of the server
	GetHostname() string

	// ServeHTTP implements the http.Handler interface, allowing the server to be used with standard HTTP libraries
	ServeHTTP(w http.ResponseWriter, r *http.Request)

	// SetConfig applies the provided configuration to the server
	SetConfig(*ServerConfig)
}

// NewServer creates a new server instance based on the provided configuration.
// It returns an implementation of the Server interface that matches the specified protocol.
// Currently supported protocols are "HTTP" and "THRIFT".
//
// The function creates the appropriate server type, applies the configuration,
// and returns the initialized server ready to have routes registered and be started.
func NewServer(config ServerConfig) Server {
	var server Server
	switch config.Protocol {
	case "THRIFT":
		server = NewThriftServer()
	case "HTTP":
		server = NewHTTPAPIServer()
	}
	server.SetConfig(&config)

	return server
}
