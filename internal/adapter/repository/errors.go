// Package repository provides adapter implementations for domain repository interfaces.
package repository

import "errors"

// Common repository errors for API operations.
// These errors are used across different repository implementations (Gmail, Calendar, etc.)
// to provide consistent error handling for common API error scenarios.
var (
	// ErrBadRequest is returned when the API request is invalid or malformed.
	ErrBadRequest = errors.New("bad request")

	// ErrRateLimited is returned when API rate limits have been exceeded.
	ErrRateLimited = errors.New("rate limited")

	// ErrTemporary is returned for temporary/transient errors that may be retried.
	ErrTemporary = errors.New("temporary error")
)
