// Package responder provides implementations for handling API responses across different protocols.
// It defines a common interface and protocol-specific implementations for HTTP and Thrift.
package responder

import "github.com/phnam/go-protocol-adapter/common"

// APIResponder defines the interface for handling API responses.
// It provides methods to format and send responses in a protocol-agnostic way,
// allowing the same business logic to work with different protocols.
type APIResponder interface {
	// Respond processes the given APIResponse and sends it to the client.
	// It handles protocol-specific formatting and transmission details.
	// Returns an error if the response cannot be processed or sent.
	Respond(*common.APIResponse[any]) error

	// GetRawResponse returns the underlying raw response object.
	// The returned interface{} can be cast to the appropriate protocol-specific type.
	GetRawResponse() interface{}

	// SetFuncName sets the function name that will be included in response headers.
	// This is useful for debugging and tracing requests through the system.
	SetFuncName(string)
}
