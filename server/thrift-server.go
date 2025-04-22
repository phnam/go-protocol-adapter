package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"

	"github.com/apache/thrift/lib/go/thrift"
	sdk "github.com/phnam/go-protocol-adapter"
	"github.com/phnam/go-protocol-adapter/common"
	requestPackage "github.com/phnam/go-protocol-adapter/request"
	responderPackage "github.com/phnam/go-protocol-adapter/responder"
	"github.com/phnam/go-protocol-adapter/thriftapi"
)

// ThriftServer implements the Server interface for the Apache Thrift protocol.
// It provides an RPC-style API interface using Thrift's binary serialization format.
type ThriftServer struct {
	// rootServer is the underlying Thrift server instance
	rootServer *thrift.TSimpleServer
	// thriftHandler handles incoming Thrift API requests
	thriftHandler *ThriftHandler
	// port is the TCP port the server listens on
	port int
	// ID is a unique identifier for this server instance
	ID int
	// hostname stores the server's hostname
	hostname string
	// config holds the server configuration
	config *ServerConfig
}

// NewThriftServer creates a new Thrift API server instance.
// It initializes the server with default configuration values and creates a ThriftHandler
// to process incoming requests. Returns an implementation of the Server interface.
func NewThriftServer() Server {
	idCounter += 1
	hostname, _ := os.Hostname()

	server := &ThriftServer{
		ID:       idCounter,
		port:     8080, // default port
		hostname: hostname,
		config: &ServerConfig{
			// Default buffer size for transport (24KB)
			BufferSize: 1024 * 24,
			// Default maximum message size (4KB)
			MessageSize: 1024 * 4,
		},
	}

	// Initialize the Thrift request handler
	server.thriftHandler = &ThriftHandler{
		Handlers: make(map[string]Handler),
		hostname: hostname,
		server:   server,
	}
	return server
}

// SetHandler registers a handler function for a specific method and path.
// For Thrift servers, the method and path are combined into a single string key
// in the format "METHOD://path" (e.g., "GET://users").
//
// Parameters:
// - method: The HTTP method (GET, POST, etc.) from common.APIMethod
// - path: The path pattern to match
// - fn: The handler function to execute when the route is matched
func (server *ThriftServer) SetHandler(method *common.MethodValue, path string, fn Handler) error {
	fullPath := string(method.Value) + "://" + path
	server.thriftHandler.Handlers[fullPath] = fn
	return nil
}

// PreRequest registers a handler function that will be executed before every request.
// This can be used for authentication, logging, or other cross-cutting concerns.
//
// The pre-request handler receives the request and can perform validation or modifications
// before the main handler is called. If the pre-request handler returns a response,
// the main handler will not be called.
func (server *ThriftServer) PreRequest(fn Handler) error {
	server.thriftHandler.preHandler = fn
	return nil
}

// Expose sets the port number that the server will listen on.
// This method must be called before Start() to configure the server's listening port.
func (server *ThriftServer) Expose(port int) {
	server.port = port
}

// Start begins listening for incoming Thrift RPC requests on the configured port.
// It sets up the Thrift server with the appropriate transport, protocol, and processor,
// then starts the server. The method blocks until the server encounters an error or is shut down.
//
// The server uses:
// - TServerSocket for the transport layer
// - TFramedTransport with buffering for framing
// - TBinaryProtocol for serialization
//
// The WaitGroup parameter allows the caller to wait for the server to exit.
// The method calls wg.Done() when the server exits, regardless of whether it
// exited due to an error or normal shutdown.
func (server *ThriftServer) Start(wg *sync.WaitGroup) {
	var ps = strconv.Itoa(server.port)
	fmt.Println("  [ Thrift Server " + strconv.Itoa(server.ID) + " ] Try to listen at " + ps)

	// Create a TCP socket transport
	var transport thrift.TServerTransport
	transport, _ = thrift.NewTServerSocket("0.0.0.0:" + ps)

	// Create a processor that will handle incoming requests
	proc := thriftapi.NewAPIServiceProcessor(server.thriftHandler)

	// Create the server with the configured transport, protocol, and processor
	server.rootServer = thrift.NewTSimpleServer4(proc, transport,
		// Use framed transport with buffering for better performance
		thrift.NewTFramedTransportFactoryConf(
			thrift.NewTBufferedTransportFactory(server.config.BufferSize),
			&thrift.TConfiguration{
				MaxFrameSize: server.config.MessageSize,
			}),
		// Use binary protocol for serialization
		thrift.NewTBinaryProtocolFactoryConf(
			&thrift.TConfiguration{
				MaxMessageSize: server.config.MessageSize,
			}))

	// Start the server (blocks until server exits)
	err := server.rootServer.Serve()
	if err != nil {
		panic(err)
	}
	wg.Done()
}

// GetHostname returns the hostname of the server.
// This is typically used for including the hostname in response headers.
func (server *ThriftServer) GetHostname() string {
	return server.hostname
}

// ServeHTTP implements the http.Handler interface for compatibility with HTTP servers.
// For Thrift servers, this method is a no-op since Thrift uses its own transport protocol.
func (server *ThriftServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	return
}

// SetConfig applies the provided configuration to the server.
// This method is called by NewServer after creating the server instance.
// It updates the server's configuration with the provided values.
func (server *ThriftServer) SetConfig(config *ServerConfig) {
	server.config = config
}

// ThriftHandler implements the Thrift service interface for handling API requests.
// It processes incoming Thrift RPC calls, maps them to the appropriate handler function,
// and returns the response in the Thrift format.
type ThriftHandler struct {
	// Handlers maps route patterns to handler functions
	Handlers map[string]Handler
	// preHandler is the optional handler function executed before every request
	preHandler Handler
	// hostname stores the server's hostname for inclusion in response headers
	hostname string
	// server is a reference to the parent Thrift server
	server *ThriftServer
}

// Call implements the Thrift service interface method for handling API requests.
// This method is called by the Thrift framework for each incoming RPC request.
//
// The method performs the following steps:
// 1. Sets up panic recovery to ensure proper error responses
// 2. Creates request and responder objects
// 3. Executes the pre-request handler if configured
// 4. Attempts to find and execute the appropriate handler for the request path
// 5. Returns the response in Thrift format
//
// If no matching handler is found, it returns a NOT_FOUND error response.
func (th *ThriftHandler) Call(ctx context.Context, request *thriftapi.APIRequest) (r *thriftapi.APIResponse, err error) {
	// Set up panic recovery to ensure we always return a proper response
	defer func() {
		if rec := recover(); rec != nil {
			r = &thriftapi.APIResponse{
				Status:    thriftapi.Status_ERROR,
				Message:   "There is an error, please try again later.",
				ErrorCode: "INTERNAL_SERVICE_ERROR",
			}

			log.Println("panic: ", rec, string(debug.Stack()))
		}
	}()

	// Create request and responder objects
	var req = requestPackage.NewThriftAPIRequest(request)
	var responder = responderPackage.NewThriftAPIResponder(th.hostname, "ThriftHandler.Call")
	var resp *thriftapi.APIResponse

	// Process pre-request handler if configured
	if th.preHandler != nil {
		// Set function name in responder for tracing/debugging
		if th.server.config == nil || !th.server.config.HideFuncName {
			responder.SetFuncName(sdk.GetFunctionName(th.preHandler))
		}

		// Execute the pre-request handler
		err := th.preHandler(req, responder)

		// Check if the pre-request handler generated a response
		tmp := responder.GetRawResponse()
		if tmp != nil {
			resp = tmp.(*thriftapi.APIResponse)
		}

		// If no response but there was an error, create an error response
		if resp == nil && err != nil {
			resp = &thriftapi.APIResponse{
				Status:  thriftapi.Status_ERROR,
				Message: "PreRequest error: " + err.Error(),
			}
		}

		// If pre-request handler generated a response, return it immediately
		if resp != nil {
			return resp, nil
		}
	}

	// Process routing - find the appropriate handler for the request
	method := req.GetMethod()
	path := request.GetPath()
	fullPath := method.Value + "://" + path

	// Check for exact match first
	if th.Handlers[fullPath] != nil {
		processFunc := th.Handlers[fullPath]

		// Set function name in responder for tracing/debugging
		funcName := ""
		if th.server.config == nil || !th.server.config.HideFuncName {
			funcName = sdk.GetFunctionName(processFunc)
		}
		responder = responderPackage.NewThriftAPIResponder(th.hostname, funcName)

		// Execute the handler
		err = processFunc(req, responder)

		// Get and return the response
		resp = nil
		tmp := responder.GetRawResponse()
		if tmp != nil {
			resp = tmp.(*thriftapi.APIResponse)
		}
		return resp, err
	} else {
		// No exact match found, try pattern matching with path parameters
		inputParts := strings.Split(path, "/")

		// Setup data for pattern matching
		var selectedHandler Handler = nil
		var selectedScore = 0
		var selectedVarCount = 0
		var varMap = map[string]string{}

		// Try to match each route pattern
		for full, hdl := range th.Handlers {
			// Initialize scoring variables for this handler
			var score = 0
			var varCount = 0
			var tempMap = map[string]string{}

			// Split the route into method and path parts
			methodPath := strings.Split(full, "://")
			// Skip if method doesn't match
			if method.Value != methodPath[0] {
				continue
			}

			// Compare each path segment
			validation := true
			pathParts := strings.Split(methodPath[1], "/")
			for i, part := range pathParts {
				if i < len(inputParts) {
					if strings.HasPrefix(part, ":") {
						// This is a path parameter
						tempMap[part[1:]] = inputParts[i]
						varCount = varCount + 1
					} else if part != inputParts[i] {
						// This segment doesn't match
						validation = false
						break
					}
					// Increment score for each matching segment
					score = i + 1 // if match at parts[0] => score = 1
				} else {
					break
				}
			}

			// Skip if validation failed
			if !validation {
				continue
			}

			// Determine if this is a better match than what we've found so far
			// Prioritize by: score, exact length match, and fewer variables
			if score > selectedScore || (score == selectedScore && len(pathParts) == len(inputParts) && varCount <= selectedVarCount) {
				varMap = tempMap
				selectedHandler = hdl
				selectedScore = score
				selectedVarCount = varCount
			}
		}

		// If we found a matching handler with pattern matching
		if selectedHandler != nil {
			// Apply URL parameters from the matched route
			for key, value := range varMap {
				req.SetVar(key, value)
			}

			// Set function name in responder for tracing/debugging
			funcName := ""
			if th.server.config == nil || !th.server.config.HideFuncName {
				funcName = sdk.GetFunctionName(selectedHandler)
			}
			responder = responderPackage.NewThriftAPIResponder(th.hostname, funcName)

			// Execute the selected handler
			err = selectedHandler(req, responder)

			// Get and return the response
			resp = nil
			tmp := responder.GetRawResponse()
			if tmp != nil {
				resp = tmp.(*thriftapi.APIResponse)
			}
			return resp, err
		}
	}

	// No matching handler found, return a NOT_FOUND error response
	return &thriftapi.APIResponse{
		Status:    thriftapi.Status_NOT_FOUND,
		Message:   "API Method/Path " + method.Value + " " + path + " isn't found",
		ErrorCode: "API_NOT_FOUND",
		Headers: map[string]string{
			"X-Hostname": th.hostname,
			"X-Function": "ThriftHandler.Call",
		},
	}, nil
}
