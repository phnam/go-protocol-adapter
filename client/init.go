// Package client provides API client implementations for different protocols.
package client

import (
	"fmt"
	"time"

	"github.com/phnam/go-protocol-adapter/common"
	sdk "github.com/phnam/go-protocol-adapter/request"
)

// APIClient defines the interface for making API requests across different protocols.
type APIClient[T any] interface {
	MakeRequest(sdk.APIRequest) *common.APIResponse[T]
	SetDebug(bool)
}

// APIClientConfiguration contains all the configuration parameters needed to create an API client.
type APIClientConfiguration struct {
	// Address is the endpoint URL or host:port of the API server
	Address string
	// Protocol specifies the communication protocol ("HTTP" or "THRIFT")
	Protocol string
	// Timeout is the maximum duration to wait for a request to complete
	Timeout time.Duration
	// MaxRetry is the maximum number of retry attempts for failed requests
	MaxRetry int
	// WaitToRetry is the duration to wait between retry attempts
	WaitToRetry time.Duration

	// MaxConnection defines the maximum number of concurrent connections (for Thrift)
	MaxConnection int
	// ErrorLogOnly when true, only logs errors and not successful requests
	ErrorLogOnly bool

	// ResultObject can hold a custom result object for the client
	ResultObject interface{}

	// KeepDataStringFormat when true, keeps response data as string format (used for Thrift client)
	KeepDataStringFormat *bool
}

// NewAPIClient creates a new API client based on the specified protocol in the configuration.
// It returns an implementation of the APIClient interface based on the protocol:
// - "THRIFT": Returns a ThriftClient
// - "HTTP": Returns a RestClient
// If an unsupported protocol is specified, it returns nil.
func NewAPIClient[T any](config *APIClientConfiguration) APIClient[T] {
	if config == nil {
		return nil
	}

	if config.Timeout < 10*time.Millisecond {
		fmt.Println("[WARNING] Timeout is too short. It should be at least 10ms.")
		config.Timeout = 10 * time.Millisecond
	}
	switch config.Protocol {
	case "THRIFT":
		return NewThriftClient[T](config)
	case "HTTP":
		return NewHTTPClient[T](config)
	}
	return nil
}
