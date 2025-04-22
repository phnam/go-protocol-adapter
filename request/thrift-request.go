package request

import (
	"encoding/json"
	"strings"

	"github.com/phnam/go-protocol-adapter/common"
	"github.com/phnam/go-protocol-adapter/thriftapi"
)

// APIThriftRequest implements the APIRequest interface for Thrift protocol.
// It wraps a thriftapi.APIRequest to provide access to Thrift request data.
type APIThriftRequest struct {
	t          string                 // Protocol type identifier
	context    *thriftapi.APIRequest  // The underlying Thrift request
	attributes map[string]interface{} // Storage for request attributes
	variables  map[string]string      // Storage for path variables
}

// NewThriftAPIRequest creates a new Thrift API request wrapper around a thriftapi.APIRequest.
// It returns an implementation of the APIRequest interface.
func NewThriftAPIRequest(e *thriftapi.APIRequest) APIRequest {
	return &APIThriftRequest{
		t:          "THRIFT",
		context:    e,
		attributes: make(map[string]interface{}),
		variables:  map[string]string{},
	}
}

// GetPath returns the request path from the Thrift context.
func (req *APIThriftRequest) GetPath() string {
	return req.context.GetPath()
}

// GetIP returns the client's IP address from the X-Forwarded-For header.
// Returns an empty string if the header is not present.
func (req *APIThriftRequest) GetIP() string {
	forwarded := req.GetHeader("X-Forwarded-For")
	if forwarded == "" {
		return ""
	}

	splitted := strings.Split(forwarded, ",")
	return splitted[0]
}

// GetMethod returns the request method as a common.MethodValue.
// It maps method strings to the application's method enum values.
func (req *APIThriftRequest) GetMethod() *common.MethodValue {
	var s = req.context.GetMethod()
	switch s {
	case "GET":
		return common.APIMethod.GET
	case "POST":
		return common.APIMethod.POST
	case "PUT":
		return common.APIMethod.PUT
	case "PATCH":
		return common.APIMethod.PATCH
	case "OPTIONS":
		return common.APIMethod.OPTIONS
	case "QUERY":
		return common.APIMethod.QUERY
	case "DELETE":
		return common.APIMethod.DELETE
	}

	return &common.MethodValue{Value: s}
}

// GetParam retrieves a query parameter by name from the request.
// Returns an empty string if the parameter doesn't exist or params are nil.
func (req *APIThriftRequest) GetParam(name string) string {
	params := req.context.GetParams()
	if params == nil {
		return ""
	}
	return params[name]
}

// GetParams returns all query parameters as a map of string keys and values.
func (req *APIThriftRequest) GetParams() map[string]string {
	return req.context.GetParams()
}

// ParseBody unmarshals the request body into the provided interface.
// It uses JSON unmarshaling to parse the request content.
func (req *APIThriftRequest) ParseBody(data interface{}) error {
	return json.Unmarshal([]byte(req.context.Content), &data)
}

// GetContentText returns the raw request body as a string.
func (req *APIThriftRequest) GetContentText() string {
	return req.context.Content
}

// GetHeader retrieves a specific header value by name.
// Returns an empty string if the header doesn't exist or headers are nil.
func (req *APIThriftRequest) GetHeader(name string) string {
	headers := req.context.GetHeaders()
	if headers == nil {
		return ""
	}
	return headers[name]
}

// GetHeaders returns all headers as a map of string keys and values.
func (req *APIThriftRequest) GetHeaders() map[string]string {
	return req.context.GetHeaders()
}

// GetAttribute retrieves a context attribute by name from the internal attributes map.
func (req *APIThriftRequest) GetAttribute(name string) interface{} {
	return req.attributes[name]
}

// SetAttribute stores a context attribute in the internal attributes map.
func (req *APIThriftRequest) SetAttribute(name string, value interface{}) {
	req.attributes[name] = value
}

// GetAttr is an alias for GetAttribute that retrieves a context attribute by name.
// This method exists for backward compatibility.
func (req *APIThriftRequest) GetAttr(name string) interface{} {
	return req.attributes[name]
}

// SetAttr is an alias for SetAttribute that stores a context attribute.
// This method exists for backward compatibility.
func (req *APIThriftRequest) SetAttr(name string, value interface{}) {
	req.attributes[name] = value
}

// GetVar retrieves a path variable/parameter by name from the internal variables map.
func (req *APIThriftRequest) GetVar(name string) string {
	return req.variables[name]
}

// SetVar sets a path variable/parameter in the internal variables map.
func (req *APIThriftRequest) SetVar(name string, value string) {
	req.variables[name] = value
}
