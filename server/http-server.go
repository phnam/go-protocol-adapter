package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	adapter "github.com/phnam/go-protocol-adapter"
	"github.com/phnam/go-protocol-adapter/common"
	"github.com/phnam/go-protocol-adapter/request"
	responderPackage "github.com/phnam/go-protocol-adapter/responder"
)

// HTTPAPIServer implements the Server interface for HTTP protocol.
// It uses the Echo framework to handle HTTP requests and provides a RESTful API interface.
type HTTPAPIServer struct {
	// T identifies the protocol type as "HTTP"
	T string
	// Echo is the underlying Echo framework instance
	Echo *echo.Echo
	// Thrift is an optional reference to a Thrift server for hybrid deployments
	Thrift *ThriftServer
	// Port is the HTTP port the server listens on
	Port int
	// ID is a unique identifier for this server instance
	ID int
	// RunSSL indicates whether HTTPS should be enabled
	RunSSL bool
	// SSLPort is the HTTPS port the server listens on when SSL is enabled
	SSLPort int
	// hostname stores the server's hostname
	hostname string
	// config holds the server configuration
	config *ServerConfig
	// debug enables verbose logging when true
	debug bool
	// router maps route patterns to handler functions
	router map[string]Handler
}

// NewHTTPAPIServer creates a new HTTP API server instance.
// It initializes the Echo framework, sets up middleware, and configures error handling.
// Returns an implementation of the Server interface.
func NewHTTPAPIServer() Server {
	idCounter += 1
	hostname, _ := os.Hostname()
	var server = HTTPAPIServer{
		T:        "HTTP",
		Echo:     echo.New(),
		ID:       idCounter,
		hostname: hostname,
		router:   map[string]Handler{},
	}
	// Enable Gzip compression for responses
	server.Echo.Use(middleware.Gzip())

	// Configure custom error handler for routes not found
	server.Echo.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}
		c.String(http.StatusNotFound, `{"status":"NOT_FOUND","error_code":"NOT_FOUND","message":"[SDK] Route not found for "`+c.Request().Method+` `+c.Request().URL.Path+`"}`)
	}
	return &server
}

// SetHandler registers a handler function for a specific HTTP method and path.
// It wraps the provided handler function in a HandlerWrapper that handles common
// functionality like error handling and response formatting.
//
// Parameters:
// - method: The HTTP method (GET, POST, etc.) from common.APIMethod
// - path: The URL path pattern to match
// - fn: The handler function to execute when the route is matched
//
// The method maps the handler to the appropriate Echo framework route and also
// stores it in the internal router map for dynamic route matching.
func (server *HTTPAPIServer) SetHandler(method *common.MethodValue, path string, fn Handler) error {
	var wrapper = &HandlerWrapper{
		handler: fn,
		server:  server,
	}

	switch method.Value {
	case common.APIMethod.GET.Value:
		server.Echo.GET(path, wrapper.processCore)
	case common.APIMethod.POST.Value:
		server.Echo.POST(path, wrapper.processCore)
	case common.APIMethod.PUT.Value:
		server.Echo.PUT(path, wrapper.processCore)
	case common.APIMethod.DELETE.Value:
		server.Echo.DELETE(path, wrapper.processCore)
	case common.APIMethod.QUERY.Value:
		server.Echo.Add(method.Value, path, wrapper.processCore)
	}
	server.router[method.Value+path] = fn

	return nil
}

// PreRequest registers a handler function that will be executed before every request.
// This can be used for authentication, logging, or other cross-cutting concerns.
//
// The method adds the handler as Echo middleware, which wraps all subsequent handlers.
// The pre-request handler receives the request and can perform validation or modifications
// before the main handler is called. If the pre-request handler returns an error,
// the main handler will not be called.
//
// The method also includes special handling for QUERY requests, attempting to find
// a matching route dynamically if the standard routing fails.
func (server *HTTPAPIServer) PreRequest(fn Handler) error {
	server.Echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			funcName := ""
			if server.config == nil || !server.config.HideFuncName {
				funcName = adapter.GetFunctionName(fn)
			}

			req := request.NewHTTPAPIRequest(c)
			responder := responderPackage.NewHTTPAPIResponder(c, server.GetHostname(), funcName)
			if server.debug {
				fmt.Println("Before PreHandlerWrapper.processCore: ", req.GetMethod(), req.GetMethod().Value, funcName)
			}

			// Set up panic recovery to ensure we always return a proper response
			defer func() {
				if server.debug {
					fmt.Println("Exit PreHandlerWrapper.processCore: ", req.GetMethod(), req.GetPath())
				}
				if r := recover(); r != nil {
					if responder != nil {
						responder.Respond(common.NewErrorResponse("ERROR", "PANIC", "Please try again later."))
					}
					log.Println("panic: ", r, string(debug.Stack()))
				}
			}()

			// Execute the pre-request handler
			err := fn(req, responder)

			if server.debug {
				fmt.Println("After PreHandlerWrapper.processCore: ", req.GetMethod().Value, err)
				fmt.Println("Next handler", next != nil)
			}

			// If pre-request handler succeeds, continue to the main handler
			if err == nil {
				err = next(c)
			}

			if server.debug {
				fmt.Println("After PreHandlerWrapper.MAIN: ", req.GetMethod().Value, err)
			}

			// Special handling for QUERY requests - try to find a matching route dynamically
			if err != nil && !c.Response().Committed && req.GetMethod().Value == "QUERY" {
				handler, varMap := findRoute(req.GetMethod().Value, req.GetPath(), server.router)
				if handler != nil {
					// Apply URL parameters from the matched route
					if varMap != nil {
						for key, value := range varMap {
							req.SetVar(key, value)
						}
					}
					if server.config == nil || !server.config.HideFuncName {
						responder.SetFuncName(adapter.GetFunctionName(handler))
					}
					handler(req, responder)
				} else {
					responder.Respond(common.NewErrorResponse("NOT_FOUND", "NOT_FOUND", "Route not found"))
				}
			}

			return nil
		}
	})
	return nil
}

// Expose sets the port number that the server will listen on for HTTP connections.
func (server *HTTPAPIServer) Expose(port int) {
	server.Port = port
}

// ExposeSSL enables HTTPS and sets the port number for SSL/TLS connections.
// This allows the server to handle both HTTP and HTTPS connections simultaneously.
func (server *HTTPAPIServer) ExposeSSL(port int) {
	server.RunSSL = true
	server.SSLPort = port
}

// Start begins listening for incoming HTTP requests on the configured port.
// If SSL is enabled, it also starts an HTTPS server on the configured SSL port.
// The method blocks until the server encounters an error or is shut down.
//
// The WaitGroup parameter allows the caller to wait for the server to exit.
// The method calls wg.Done() when the server exits, regardless of whether it
// exited due to an error or normal shutdown.
func (server *HTTPAPIServer) Start(wg *sync.WaitGroup) {
	var ps = strconv.Itoa(server.Port)
	fmt.Println("  [ API Server " + strconv.Itoa(server.ID) + " ] Try to listen at " + ps)
	server.Echo.HideBanner = true

	// Start HTTPS server in a separate goroutine if SSL is enabled
	if server.RunSSL {
		go func() {
			err := server.Echo.StartTLS(":"+strconv.Itoa(server.SSLPort), "crt.pem", "key.pem")
			if err != nil {
				fmt.Println("[Start TLS error] " + err.Error())
			}
		}()
	}

	// Start HTTP server (blocks until server exits)
	err := server.Echo.Start(":" + ps)
	if err != nil {
		fmt.Println("Fail to start " + err.Error())
	}
	wg.Done()
}

// GetHostname returns the hostname of the server.
// This is typically used for including the hostname in response headers.
func (server *HTTPAPIServer) GetHostname() string {
	return server.hostname
}

// ServeHTTP implements the http.Handler interface, allowing the server to be used with standard HTTP libraries.
// It delegates to the underlying Echo framework's ServeHTTP method.
func (server *HTTPAPIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.Echo.ServeHTTP(w, r)
}

// SetConfig applies the provided configuration to the server.
// This method is called by NewServer after creating the server instance.
func (server *HTTPAPIServer) SetConfig(config *ServerConfig) {
	server.config = config
}

// HandlerWrapper wraps a handler function with common functionality like error handling.
// It adapts between the Echo framework's handler interface and the application's handler interface.
type HandlerWrapper struct {
	// handler is the application-specific handler function to execute
	handler Handler
	// server is a reference to the parent HTTP server
	server *HTTPAPIServer
}

// Handler defines the function signature for API request handlers.
// It takes an APIRequest and APIResponder and returns an error if the handling fails.
// This is the core handler type used throughout the application for processing requests.
type Handler = func(req request.APIRequest, res responderPackage.APIResponder) error

// processCore is the Echo framework handler function that wraps the application handler.
// It creates the appropriate request and responder objects, calls the handler,
// and handles any panics that might occur during processing.
func (hw *HandlerWrapper) processCore(c echo.Context) error {
	if hw.server.debug {
		fmt.Println("Start MAIN.processCore: ", c.Request().Method, c.Request().URL.Path)
	}

	// Get the function name for debugging/tracing if not disabled
	funcName := ""
	if hw.server.config == nil || !hw.server.config.HideFuncName {
		funcName = adapter.GetFunctionName(hw.handler)
	}

	// Create request and responder objects
	req := request.NewHTTPAPIRequest(c)
	responder := responderPackage.NewHTTPAPIResponder(c, hw.server.GetHostname(), funcName)

	if hw.server.debug {
		fmt.Println("Before MAIN.processCore: ", req.GetMethod(), req.GetMethod().Value, funcName)
	}

	// Set up panic recovery to ensure we always return a proper response
	defer func() {
		if r := recover(); r != nil {
			if responder != nil {
				responder.Respond(common.NewErrorResponse("ERROR", "PANIC", "Please try again later."))
			}
			log.Println("panic: ", r, string(debug.Stack()))
		}
	}()

	// Execute the handler
	hw.handler(req, responder)

	if hw.server.debug {
		fmt.Println("After MAIN.processCore: ", req.GetMethod(), req.GetMethod().Value, funcName)
	}

	return nil
}

// PreHandlerWrapper wraps a pre-request handler with Echo middleware functionality.
// This type is used internally by the PreRequest method to adapt between
// the Echo middleware interface and the application's handler interface.
type PreHandlerWrapper struct {
	// preHandler is the application-specific pre-request handler function
	preHandler Handler
	// next is the next Echo handler in the middleware chain
	next echo.HandlerFunc
	// server is a reference to the parent HTTP server
	server *HTTPAPIServer
	// funcName is the name of the handler function for debugging/tracing
	funcName string
}

// SetDebug enables or disables debug logging for the server.
// When debug is enabled, the server will print detailed information about request processing.
func (server *HTTPAPIServer) SetDebug(debug bool) {
	server.debug = debug
}

// findRoute attempts to find a matching route handler for the given method and path.
// It supports path parameters (e.g., "/users/:id") and returns both the handler
// and a map of parameter names to values.
//
// The function first checks for an exact match. If none is found, it tries to match
// routes with path parameters, using a scoring system to find the best match:
//  1. Routes with more matching segments have higher priority
//  2. For routes with the same number of matching segments, those with fewer variables are preferred
//  3. For routes with the same number of matching segments and variables, those with variables
//     appearing later in the path are preferred
//
// Returns the matched handler and a map of path parameters, or nil if no match is found.
func findRoute(method string, path string, handlerMap map[string]Handler) (Handler, map[string]string) {
	if handlerMap == nil {
		return nil, nil
	}
	// Check for exact match first
	if handlerMap[method+path] != nil {
		return handlerMap[method+path], nil
	}

	// Prepare for pattern matching
	targetRoute := method + path
	targetParts := strings.Split(targetRoute, "/")
	var selectedHandler Handler
	currentScore := 0
	var currentVarMap map[string]string
	currentFirstVar := 0

	// Try to match each route pattern
	for route, handler := range handlerMap {
		varMap := map[string]string{}
		parts := strings.Split(route, "/")
		score := 0
		firstVar := 0

		// Compare each path segment
		for i, part := range parts {
			if part[0] == ':' {
				// This is a path parameter
				varMap[part[1:]] = targetParts[i]
				if firstVar == 0 {
					firstVar = i
				}
			} else if part != targetParts[i] {
				// This segment doesn't match
				break
			}
			score++
		}

		// If we didn't match all segments of the route, skip it
		if score < len(parts) {
			continue
		}

		// Determine if this is a better match than what we've found so far
		if score > currentScore || (score == currentScore && len(varMap) < len(currentVarMap)) ||
			(score == currentScore && len(varMap) == len(currentVarMap) && firstVar > currentFirstVar) {
			selectedHandler = handler
			currentScore = score
			currentVarMap = varMap
			currentFirstVar = firstVar
		}
	}

	if selectedHandler != nil {
		return selectedHandler, currentVarMap
	}

	return nil, nil
}
