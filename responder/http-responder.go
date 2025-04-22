package responder

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/labstack/echo"
	"github.com/phnam/go-protocol-adapter/common"
)

// HTTPAPIResponder implements the APIResponder interface for HTTP protocol.
// It handles formatting and sending API responses over HTTP using the Echo framework.
type HTTPAPIResponder struct {
	// t identifies the protocol type as "HTTP"
	t string
	// context is the Echo context for the current request
	context echo.Context
	// start tracks when the request processing began for calculating execution time
	start time.Time
	// hostname stores the server hostname to include in response headers
	hostname string
	// funcName stores the handler function name to include in response headers
	funcName string
	// resp stores the raw response object after it's been sent
	resp interface{}
}

// NewHTTPAPIResponder creates a new HTTP API responder with the given Echo context, hostname, and function name.
// It initializes a timer to track execution time and returns an implementation of the APIResponder interface.
func NewHTTPAPIResponder(c echo.Context, hostname string, funcName string) APIResponder {
	return &HTTPAPIResponder{
		t:        "HTTP",
		start:    time.Now(),
		context:  c,
		hostname: hostname,
		funcName: funcName,
	}
}

// Respond processes and sends the API response to the client over HTTP.
// It validates the response, sets appropriate headers, maps API status to HTTP status codes,
// and sends the response as JSON (or redirects for redirected status).
//
// The method performs the following steps:
// 1. Validates that the response is not nil and data is a slice
// 2. Copies any headers from the response to the HTTP response
// 3. Adds execution time, hostname, and function name headers
// 4. Maps the API status to the appropriate HTTP status code
// 5. Sends the response with the correct content type
//
// Returns an error if the response cannot be processed or sent.
func (resp *HTTPAPIResponder) Respond(response *common.APIResponse[any]) error {
	var context = resp.context

	if response == nil {
		return errors.New("response cannot be nil")
	}

	if response.Data != nil && reflect.TypeOf(response.Data).Kind() != reflect.Slice {
		return errors.New("data response must be a slice")
	}

	if response.Headers != nil {
		header := context.Response().Header()
		for key, value := range response.Headers {
			header.Set(key, value)
		}
		response.Headers = nil
	}

	var dif = float64(time.Since(resp.start).Nanoseconds()) / 1000000
	context.Response().Header().Set("X-Execution-Time", fmt.Sprintf("%.4f ms", dif))
	context.Response().Header().Set("X-Hostname", resp.hostname)

	if resp.funcName != "" {
		context.Response().Header().Set("X-Function", resp.funcName)
	}

	switch response.Status {
	case common.APIStatus.Ok:
		return context.JSON(http.StatusOK, response)
	case common.APIStatus.Error:
		return context.JSON(http.StatusInternalServerError, response)
	case common.APIStatus.Forbidden:
		return context.JSON(http.StatusForbidden, response)
	case common.APIStatus.Invalid:
		return context.JSON(http.StatusBadRequest, response)
	case common.APIStatus.NotFound:
		return context.JSON(http.StatusNotFound, response)
	case common.APIStatus.Unauthorized:
		return context.JSON(http.StatusUnauthorized, response)
	case common.APIStatus.Existed:
		return context.JSON(http.StatusConflict, response)
	case common.APIStatus.Redirected:
		return context.Redirect(http.StatusFound, context.Response().Header().Get("Location"))
	}

	resp.resp = response

	return context.JSON(http.StatusBadRequest, response)
}

// GetRawResponse returns the underlying raw response object.
// This can be used to access the response after it has been sent.
func (resp *HTTPAPIResponder) GetRawResponse() interface{} {
	return resp.resp
}

// SetFuncName sets the function name that will be included in the X-Function response header.
// This is useful for debugging and tracing requests through the system.
func (resp *HTTPAPIResponder) SetFuncName(name string) {
	resp.funcName = name
}
