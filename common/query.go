// Package common provides shared types, constants, and utilities used across the protocol adapter.
package common

// Query represents a generic query structure for data retrieval operations.
// It provides a flexible way to filter, paginate, sort, and customize data queries.
// The generic type parameter T allows for type-safe filter definitions.
type Query[T any] struct {
	Filter  T                      `json:"filter,omitempty"`     // Type-specific filter criteria
	Offset  int64                  `json:"offset,omitempty"`     // Number of records to skip for pagination
	Limit   int64                  `json:"limit,omitempty"`      // Maximum number of records to return
	Sort    map[string]int         `json:"sort_field,omitempty"` // Field-based sorting (field name -> direction)
	Options map[string]interface{} `json:"options,omitempty"`    // Additional query options
}
