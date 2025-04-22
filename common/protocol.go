// Package common provides shared types, constants, and utilities used across the protocol adapter.
package common

// ProtocolEnum defines a structure containing supported communication protocols.
// This enum is used to specify which protocol to use for client-server communication.
type ProtocolEnum struct {
	HTTP   string // HTTP protocol identifier
	THRIFT string // Apache Thrift protocol identifier
}

// Protocol is a published enum containing predefined protocol values.
// It provides a convenient way to reference supported protocols throughout the application.
var Protocol = ProtocolEnum{
	HTTP:   "HTTP",
	THRIFT: "THRIFT",
}
