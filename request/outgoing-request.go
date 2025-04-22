package request

import (
	"encoding/json"

	"github.com/phnam/go-protocol-adapter/common"
)

// OutboundAPIRequest implements the APIRequest interface for outgoing requests to other services.
// It contains all the necessary data to make an API call to another service.
type OutboundAPIRequest struct {
	Method  string            `json:"method" bson:"method"`                       // HTTP method or operation type
	Path    string            `json:"path" bson:"path"`                           // Request path/endpoint
	Params  map[string]string `json:"params,omitempty" bson:"params,omitempty"`   // Query parameters
	Headers map[string]string `json:"headers,headers" bson:"headers,omitempty"`   // HTTP headers
	Content string            `json:"content,omitempty" bson:"content,omitempty"` // Request body content
}

// NewOutboundAPIRequest creates a new outbound API request with the specified parameters.
// It returns an implementation of the APIRequest interface for making calls to other services.
func NewOutboundAPIRequest(method string, path string, params map[string]string, content string, headers map[string]string) APIRequest {
	return &OutboundAPIRequest{
		Method:  method,
		Path:    path,
		Params:  params,
		Content: content,
		Headers: headers,
	}
}

// GetPath returns the request path/endpoint.
func (req *OutboundAPIRequest) GetPath() string {
	return req.Path
}

// GetIP returns a placeholder string as IP is not applicable for outbound requests.
func (req *OutboundAPIRequest) GetIP() string {
	return "GetIP() not implemented"
}

// GetMethod returns the request method as a common.MethodValue.
// It maps method strings to the application's method enum values.
func (req *OutboundAPIRequest) GetMethod() *common.MethodValue {
	var s = req.Method
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

// GetVar retrieves a path variable/parameter by name from the params map.
// For outbound requests, path variables are stored in the same map as query parameters.
func (req *OutboundAPIRequest) GetVar(name string) string {
	return req.Params[name]
}

// GetParam retrieves a query parameter by name from the params map.
func (req *OutboundAPIRequest) GetParam(name string) string {
	return req.Params[name]
}

// GetParams returns all query parameters as a map of string keys and values.
func (req *OutboundAPIRequest) GetParams() map[string]string {
	return req.Params
}

// ParseBody unmarshals the request body into the provided interface.
// It uses JSON unmarshaling to parse the request content.
func (req *OutboundAPIRequest) ParseBody(data interface{}) error {
	json.Unmarshal([]byte(req.Content), &data)
	return nil
}

// GetContentText returns the raw request body as a string.
func (req *OutboundAPIRequest) GetContentText() string {
	return req.Content
}

// GetHeader retrieves a specific header value by name.
func (req *OutboundAPIRequest) GetHeader(name string) string {
	return req.Headers[name]
}

// GetHeaders returns all headers as a map of string keys and values.
func (req *OutboundAPIRequest) GetHeaders() map[string]string {
	return req.Headers
}

// GetAttribute retrieves a context attribute by name.
// This is a no-op for outbound requests and always returns nil.
func (req *OutboundAPIRequest) GetAttribute(name string) interface{} {
	return nil
}

// SetAttribute stores a context attribute.
// This is a no-op for outbound requests as attributes are not supported.
func (req *OutboundAPIRequest) SetAttribute(name string, value interface{}) {

}

// GetAttr is an alias for GetAttribute that retrieves a context attribute by name.
// This method exists for backward compatibility and always returns nil for outbound requests.
func (req *OutboundAPIRequest) GetAttr(name string) interface{} {
	return nil
}

// SetAttr is an alias for SetAttribute that stores a context attribute.
// This method exists for backward compatibility and is a no-op for outbound requests.
func (req *OutboundAPIRequest) SetAttr(name string, value interface{}) {
	// do nothing
}

// SetVar sets a path variable/parameter.
// This is a no-op for outbound requests as variables cannot be modified after creation.
func (req *OutboundAPIRequest) SetVar(name string, value string) {
	// do nothing
}
