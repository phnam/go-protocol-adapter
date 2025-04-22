package request

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/labstack/echo"
	"github.com/phnam/go-protocol-adapter/common"
)

// HTTPAPIRequest implements the APIRequest interface for HTTP protocol.
// It wraps an echo.Context to provide access to HTTP request data.
type HTTPAPIRequest struct {
	t       string       // Protocol type identifier
	context echo.Context // The underlying Echo framework context
	body    string       // Cached request body content
}

// NewHTTPAPIRequest creates a new HTTP API request wrapper around an echo.Context.
// It returns an implementation of the APIRequest interface.
func NewHTTPAPIRequest(e echo.Context) APIRequest {
	return &HTTPAPIRequest{
		t:       "HTTP",
		context: e,
	}
}

// GetPath returns the request path from the Echo context.
func (req *HTTPAPIRequest) GetPath() string {
	return req.context.Path()
}

// GetMethod returns the HTTP method as a common.MethodValue.
// It maps standard HTTP methods to the application's method enum values.
func (req *HTTPAPIRequest) GetMethod() *common.MethodValue {
	var s = req.context.Request().Method
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
	case "DELETE":
		return common.APIMethod.DELETE
	}

	return &common.MethodValue{Value: s}
}

// GetVar retrieves a path parameter by name from the Echo context.
func (req *HTTPAPIRequest) GetVar(name string) string {
	return req.context.Param(name)
}

// SetVar sets a path parameter value.
// Note: This is a no-op in the HTTP implementation as Echo doesn't support setting path parameters.
func (req *HTTPAPIRequest) SetVar(name string, value string) {

}

// GetParam retrieves a query parameter by name from the request URL.
func (req *HTTPAPIRequest) GetParam(name string) string {
	return req.context.QueryParam(name)
}

// GetParams returns all query parameters as a map of string keys and values.
// It converts the Echo QueryParams to a simple string map.
func (req *HTTPAPIRequest) GetParams() map[string]string {
	var vals = req.context.QueryParams()
	var m = make(map[string]string)
	for key := range vals {
		m[key] = vals.Get(key)
	}
	return m
}

// ParseBody unmarshals the request body into the provided interface.
// It uses JSON unmarshaling to parse the request body content.
func (req *HTTPAPIRequest) ParseBody(data interface{}) error {

	return json.Unmarshal([]byte(req.GetContentText()), data)
}

// GetContentText returns the raw request body as a string.
// It lazily loads and caches the body content on first access.
func (req *HTTPAPIRequest) GetContentText() string {
	if req.body == "" {
		var bodyBytes []byte
		if req.context.Request().Body != nil {
			bodyBytes, _ = io.ReadAll(req.context.Request().Body)
		}

		req.body = string(bodyBytes)
	}

	return req.body
}

// GetHeader retrieves a specific HTTP header value by name.
func (req *HTTPAPIRequest) GetHeader(name string) string {
	return req.context.Request().Header.Get(name)
}

// GetHeaders returns all HTTP headers as a map of string keys and values.
// It converts the HTTP header structure to a simple string map.
func (req *HTTPAPIRequest) GetHeaders() map[string]string {
	var vals = req.context.Request().Header
	var m = make(map[string]string)
	for key := range vals {
		m[key] = vals.Get(key)
	}
	return m
}

// GetAttribute retrieves a context attribute by name from the Echo context.
func (req *HTTPAPIRequest) GetAttribute(name string) interface{} {
	return req.context.Get(name)
}

// SetAttribute stores a context attribute in the Echo context.
func (req *HTTPAPIRequest) SetAttribute(name string, value interface{}) {
	req.context.Set(name, value)
}

// GetIP returns the client's IP address.
// It first checks for X-Forwarded-For header (for proxied requests),
// then falls back to the remote address from the request.
func (req *HTTPAPIRequest) GetIP() string {
	// for forwarded case
	forwarded := req.GetHeader("X-Forwarded-For")
	if forwarded == "" {
		httpReq := req.context.Request()
		return strings.Split(httpReq.RemoteAddr, ":")[0]
	}

	splitted := strings.Split(forwarded, ",")
	return splitted[0]
}
