# Go Protocol Adapter

A flexible and high-performance protocol adapter library for Go applications that provides a unified API for HTTP and Apache Thrift protocols.

## Overview

Go Protocol Adapter is a library that abstracts away the differences between HTTP and Thrift protocols, allowing developers to write protocol-agnostic code. It provides a consistent interface for handling requests and responses, making it easy to build services that can communicate using either protocol without changing the application logic.

## Features

- **Protocol Abstraction**: Unified API for both HTTP and Thrift protocols
- **Server Implementations**: Ready-to-use server implementations for both protocols
- **Client Implementations**: Protocol-specific client implementations with a common interface
- **Request/Response Handling**: Consistent request and response handling across protocols
- **Middleware Support**: Pre-request handlers for cross-cutting concerns like authentication
- **Error Handling**: Centralized error handling with consistent error responses
- **SSL Support**: Built-in SSL/TLS support for secure communications
- **Configurable**: Extensive configuration options for both server and client

## Installation

```bash
go get github.com/phnam/go-protocol-adapter
```

## Quick Start

### Creating a Server

```go
package main

import (
    "sync"
    
    "github.com/phnam/go-protocol-adapter/common"
    "github.com/phnam/go-protocol-adapter/request"
    "github.com/phnam/go-protocol-adapter/responder"
    "github.com/phnam/go-protocol-adapter/server"
)

func main() {
    // Create a new HTTP server
    httpServer := server.NewServer(server.ServerConfig{
        Protocol: common.Protocol.HTTP,
    })
    
    // Register a handler for GET requests to the root path
    httpServer.SetHandler(common.APIMethod.GET, "/", func(req request.APIRequest, res responder.APIResponder) error {
        return res.Respond(&common.APIResponse[any]{
            Status:  common.APIStatus.Ok,
            Message: "Hello world",
            Data:    []any{map[string]interface{}{"message": "Hello from HTTP Server"}},
        })
    })
    
    // Set the port and start the server
    httpServer.Expose(8080)
    
    // Use a WaitGroup to keep the main function from exiting
    var wg sync.WaitGroup
    wg.Add(1)
    go httpServer.Start(&wg)
    
    wg.Wait()
}
```

### Creating a Client

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/phnam/go-protocol-adapter/client"
    "github.com/phnam/go-protocol-adapter/common"
    "github.com/phnam/go-protocol-adapter/request"
)

type ResponseData struct {
    Message string `json:"message"`
}

func main() {
    // Create a new HTTP client
    httpClient := client.NewAPIClient[ResponseData](&client.APIClientConfiguration{
        Address:      "localhost:8080",
        Timeout:      100 * time.Millisecond,
        MaxRetry:     1,
        Protocol:     common.Protocol.HTTP,
        ErrorLogOnly: false,
    })
    
    // Make a request
    resp := httpClient.MakeRequest(&request.OutboundAPIRequest{
        Method: "GET",
        Path:   "/",
    })
    
    // Handle the response
    if resp.Status == common.APIStatus.Ok {
        fmt.Println("Response:", resp.Data[0].Message)
    } else {
        fmt.Println("Error:", resp.Message)
    }
}
```

## Switching Protocols

One of the key benefits of this library is the ability to switch between protocols with minimal code changes. To switch from HTTP to Thrift (or vice versa), simply change the protocol in the server and client configuration:

```go
// HTTP Server
httpServer := server.NewServer(server.ServerConfig{
    Protocol: common.Protocol.HTTP,
})

// Thrift Server
thriftServer := server.NewServer(server.ServerConfig{
    Protocol: common.Protocol.THRIFT,
})

// HTTP Client
httpClient := client.NewAPIClient[ResponseData](&client.APIClientConfiguration{
    Protocol: common.Protocol.HTTP,
    // other configuration...
})

// Thrift Client
thriftClient := client.NewAPIClient[ResponseData](&client.APIClientConfiguration{
    Protocol: common.Protocol.THRIFT,
    // other configuration...
})
```

## Server Configuration

The `ServerConfig` struct provides various configuration options for servers:

```go
type ServerConfig struct {
    // Protocol specifies which protocol implementation to use ("HTTP" or "THRIFT")
    Protocol string
    
    // HideFuncName determines whether function names should be included in response headers
    HideFuncName bool
    
    // BufferSize specifies the buffer size in bytes for Thrift server transport
    BufferSize int
    
    // MessageSize specifies the maximum message size in bytes for Thrift server
    MessageSize int32
}
```

## Client Configuration

The `APIClientConfiguration` struct provides configuration options for clients:

```go
type APIClientConfiguration struct {
    // Address is the server address to connect to
    Address string
    
    // Timeout is the request timeout duration
    Timeout time.Duration
    
    // MaxRetry is the maximum number of retry attempts for failed requests
    MaxRetry int
    
    // Protocol specifies which protocol to use ("HTTP" or "THRIFT")
    Protocol string
    
    // Other configuration options...
}
```

## Examples

See the `test` directory for examples of how to use the library.
