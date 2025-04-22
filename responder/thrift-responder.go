package responder

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/phnam/go-protocol-adapter/common"
	"github.com/phnam/go-protocol-adapter/thriftapi"
)

// ThriftAPIResponder implements the APIResponder interface for the Thrift protocol.
// It handles formatting and sending API responses over Thrift, converting between
// the common APIResponse format and the Thrift-specific APIResponse format.
type ThriftAPIResponder struct {
	// t identifies the protocol type as "THRIFT"
	t string
	// resp holds the Thrift-specific response object
	resp *thriftapi.APIResponse
	// start tracks when the request processing began for calculating execution time
	start time.Time
	// hostname stores the server hostname to include in response headers
	hostname string
	// funcName stores the handler function name to include in response headers
	funcName string
}

// NewThriftAPIResponder creates a new Thrift API responder with the given hostname and function name.
// It initializes a timer to track execution time and returns an implementation of the APIResponder interface.
func NewThriftAPIResponder(hostname string, funcName string) APIResponder {
	return &ThriftAPIResponder{
		t:        "THRIFT",
		start:    time.Now(),
		hostname: hostname,
		funcName: funcName,
	}
}

// GetRawResponse returns the underlying Thrift APIResponse object.
// This can be used to access the response after it has been created.
func (responder *ThriftAPIResponder) GetRawResponse() interface{} {
	return responder.resp
}

// Respond processes the common APIResponse and converts it to a Thrift-specific APIResponse.
// It validates the response, converts the data to JSON, sets appropriate headers,
// and prepares the response for transmission over Thrift.
//
// The method performs the following steps:
// 1. Validates that the response is not nil and data is a slice
// 2. Creates a new Thrift APIResponse with the common response's fields
// 3. Converts the common status to a Thrift status enum value
// 4. Serializes the data to JSON and stores it as a string in the Content field
// 5. Adds execution time, hostname, and function name headers
//
// Returns an error if the response cannot be processed.
func (responder *ThriftAPIResponder) Respond(response *common.APIResponse[any]) error {

	if response == nil {
		return errors.New("response cannot be nil")
	}

	if response.Data != nil && reflect.TypeOf(response.Data).Kind() != reflect.Slice {
		return errors.New("data response must be a slice")
	}

	var dif = float64(time.Since(responder.start).Nanoseconds()) / 1000000

	responder.resp = &thriftapi.APIResponse{
		ErrorCode: response.ErrorCode,
		Message:   response.Message,
		Total:     response.Total,
		Headers:   response.Headers,
	}
	responder.resp.Status, _ = thriftapi.StatusFromString(response.Status)
	bytes, _ := json.Marshal(response.Data)
	responder.resp.Content = string(bytes)
	if responder.resp.Headers == nil {
		responder.resp.Headers = make(map[string]string)
	}
	responder.resp.Headers["X-Execution-Time"] = fmt.Sprintf("%.4f ms", dif)
	responder.resp.Headers["X-Hostname"] = responder.hostname

	if responder.funcName != "" {
		responder.resp.Headers["X-Function"] = responder.funcName
	}

	return nil
}

// SetFuncName sets the function name that will be included in the X-Function response header.
// This is useful for debugging and tracing requests through the system.
func (responder *ThriftAPIResponder) SetFuncName(funcName string) {
	responder.funcName = funcName
}
