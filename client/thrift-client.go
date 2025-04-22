// Package client provides API client implementations for different protocols.
package client

import (
	"context"
	"encoding/json"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/phnam/go-protocol-adapter/common"
	sdk "github.com/phnam/go-protocol-adapter/request"
	"github.com/phnam/go-protocol-adapter/thriftapi"
)

// ThriftClient implements the APIClient interface for Thrift protocol communication.
type ThriftClient[T any] struct {
	// adr is the address of the Thrift server in host:port format
	adr string
	// timeout is the maximum duration to wait for a request to complete
	timeout time.Duration
	// maxConnection is the maximum number of concurrent connections to maintain
	maxConnection int
	// maxRetry is the maximum number of retry attempts for failed requests
	maxRetry int
	// waitToRetry is the duration to wait between retry attempts
	waitToRetry time.Duration
	// cons is a map of connection IDs to ThriftCon objects
	cons map[string]*ThriftCon
	// debug enables debug logging when true
	debug bool
	// lock is a mutex for thread-safe access to the connections map
	lock *sync.Mutex
	// maxAge is the maximum age of a connection in seconds before it's refreshed
	maxAge int
	// skipUnmarshal when true, keeps response data as string format
	skipUnmarshal bool

	config *APIClientConfiguration
}

// ThriftCon represents a single connection to a Thrift API server.
type ThriftCon struct {
	// Client is the Thrift API service client
	Client *thriftapi.APIServiceClient
	// socket is the underlying transport for the connection
	socket *thrift.TTransport
	// inUsed indicates whether the connection is currently being used
	inUsed bool
	// hasError indicates whether the connection has encountered an error
	hasError bool
	// lock is a mutex for thread-safe access to this connection
	lock *sync.Mutex
	// id is the unique identifier for this connection
	id string
	// createdTime is when this connection was created
	createdTime time.Time
}

// NewThriftClient creates a new Thrift client based on the provided configuration.
// It implements the APIClient interface for Thrift protocol communication.
//
// Parameters:
//   - config: Configuration parameters for the Thrift client
//
// Returns:
//   - A pointer to a new ThriftClient instance
func NewThriftClient[T any](config *APIClientConfiguration) *ThriftClient[T] {
	// Determine whether to skip unmarshaling based on configuration
	skipUnmarshal := false
	if config.KeepDataStringFormat != nil {
		skipUnmarshal = *config.KeepDataStringFormat
	}

	// Create and return a new ThriftClient with the provided configuration
	return &ThriftClient[T]{
		adr:           config.Address,
		timeout:       config.Timeout,
		maxConnection: config.MaxConnection,
		maxRetry:      config.MaxRetry,
		waitToRetry:   config.WaitToRetry,
		cons:          make(map[string]*ThriftCon),
		lock:          &sync.Mutex{},
		maxAge:        600, // Default max age of 10 minutes
		skipUnmarshal: skipUnmarshal,
	}
}

// SetDebug enables or disables debug logging for the ThriftClient.
//
// Parameters:
//   - val: true to enable debug logging, false to disable
func (client *ThriftClient[T]) SetDebug(val bool) {
	client.debug = val
}

// newThriftCon creates a new Thrift connection to the server.
//
// Returns:
//   - A pointer to a new ThriftCon instance
func (client *ThriftClient[T]) newThriftCon() *ThriftCon {
	// Create a binary protocol factory
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	// Resolve the server address
	addr, _ := net.ResolveTCPAddr("tcp", client.adr)

	// Create a socket transport with timeout configuration
	var transport thrift.TTransport
	transport = thrift.NewTSocketFromAddrConf(addr, &thrift.TConfiguration{
		ConnectTimeout: client.timeout,
		SocketTimeout:  client.timeout,
	},
	)

	// Create a framed transport with buffering
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTBufferedTransportFactory(8192))
	transport, _ = transportFactory.GetTransport(transport)

	// Get input and output protocols
	iprot := protocolFactory.GetProtocol(transport)
	oprot := protocolFactory.GetProtocol(transport)

	// Open the transport connection
	transport.Open()

	// Create and return a new ThriftCon
	return &ThriftCon{
		socket:      &transport,
		Client:      thriftapi.NewAPIServiceClient(thrift.NewTStandardClient(iprot, oprot)),
		inUsed:      false,
		lock:        &sync.Mutex{},
		hasError:    false,
		createdTime: time.Now(),
	}
}

// pickCon selects an available connection from the pool or creates a new one.
//
// Parameters:
//   - useOld: When true, tries to reuse an existing connection before creating a new one
//
// Returns:
//   - A pointer to a ThriftCon that is ready to use, or nil if no connection could be obtained
func (client *ThriftClient[T]) pickCon(useOld bool) *ThriftCon {
	if useOld {
		client.lock.Lock()
		for conID, con := range client.cons {
			// verify if connection is free
			con.lock.Lock()
			if (*con.socket).IsOpen() {
				if !con.inUsed {
					con.inUsed = true
					con.lock.Unlock()
					client.lock.Unlock()
					return con
				}
			} else {
				delete(client.cons, conID)
				(*con.socket).Close()
			}
			con.lock.Unlock()
		}
		if len(client.cons) < client.maxConnection || client.maxConnection == 0 {
			useOld = false
		}

		client.lock.Unlock()
	}

	if !useOld {

		// if not find any available connection, create new
		con := client.newThriftCon()
		con.inUsed = true

		// append to connection pool if have space
		if len(client.cons) < client.maxConnection {
			id := rand.Intn(999999999) + 1000000000
			for client.cons[strconv.Itoa(id)] != nil {
				id = rand.Intn(999999999) + 1000000000
			}
			con.id = strconv.Itoa(id)
			client.lock.Lock()
			client.cons[con.id] = con
			client.lock.Unlock()
		}

		return con
	}

	return nil
}

// call makes a Thrift API call with the given request.
// It handles connection management and error handling.
//
// Parameters:
//   - req: The API request to process
//   - useNewCon: When true, forces the use of a new connection
//
// Returns:
//   - A pointer to a thriftapi.APIResponse containing the response
//   - An error if the call fails
func (client *ThriftClient[T]) call(req sdk.APIRequest, useNewCon bool) (*thriftapi.APIResponse, error) {

	// map to thrift request
	var r = &thriftapi.APIRequest{
		Path:    req.GetPath(),
		Params:  req.GetParams(),
		Headers: req.GetHeaders(),
		Method:  req.GetMethod().Value,
	}

	if r.Method != "GET" && r.Method != "DELETE" {
		r.Content = req.GetContentText()
	}

	// pick available connection
	var con *ThriftCon
	con = client.pickCon(!useNewCon)
	var retryToGetCon = 0
	for retryToGetCon < 10 && con == nil {
		time.Sleep(10 * time.Millisecond)
		con = client.pickCon(!useNewCon)
		retryToGetCon++
	}

	if con == nil {
		return &thriftapi.APIResponse{
			Status:  500,
			Message: "Connection pool is temporary overloaded!",
		}, &common.Error{ErrorCode: "OVERLOAD", Message: "Connection pool is overloaded! Fail to make request to " + req.GetPath()}
	}
	result, err := con.Client.Call(context.Background(), r)

	// verify error
	if err == nil {
		if con.createdTime.Add(time.Duration(client.maxAge) * time.Second).Before(time.Now()) {
			// if too old, replace this con by new con
			client.lock.Lock()
			(*con.socket).Close()
			id := con.id
			con = client.newThriftCon()
			client.cons[id] = con
			client.lock.Unlock()
		}
		con.inUsed = false
	} else {

		// remove connection from pool
		con.hasError = true
		client.lock.Lock()
		(*con.socket).Close()
		delete(client.cons, con.id)
		client.lock.Unlock()
	}

	return result, err
}

// MakeRequest implements the APIClient interface method for making API requests.
// It handles retries and error handling for Thrift service calls.
//
// Parameters:
//   - req: The API request to process
//
// Returns:
//   - A pointer to a common.APIResponse containing the response
func (client *ThriftClient[T]) MakeRequest(req sdk.APIRequest) *common.APIResponse[T] {
	now := time.Now()
	canRetry := client.maxRetry
	result, err := client.call(req, false)

	// free retry immediately if connection is not open or last connection was failed
	if err != nil {

		errMsg := strings.ToLower(err.Error())
		if (strings.Contains(errMsg, "connection not open") || strings.Contains(errMsg, "eof") ||
			strings.Contains(errMsg, "connection timed out") || strings.Contains(errMsg, "i/o timeout") ||
			strings.HasPrefix(errMsg, "overload") || strings.Contains(errMsg, "broken pipe")) && time.Now().Before(now.Add(10*time.Millisecond)) {
			result, err = client.call(req, true)
		}
	}

	// retry if failed
	for err != nil && canRetry > 0 {
		time.Sleep(client.waitToRetry)
		canRetry--
		result, err = client.call(req, true)
	}

	if err != nil {
		return &common.APIResponse[T]{
			Status:  common.APIStatus.Error,
			Message: "Endpoint error: " + err.Error(),
		}
	}

	// parse result
	resp := &common.APIResponse[T]{
		Status:    result.GetStatus().String(),
		Message:   result.GetMessage(),
		Headers:   result.GetHeaders(),
		Total:     result.GetTotal(),
		ErrorCode: result.GetErrorCode(),
		Data:      []T{},
	}
	json.Unmarshal([]byte(result.GetContent()), &resp.Data)
	return resp
}
