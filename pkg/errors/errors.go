// Package errors provides common error types for the application.
package errors

import "errors"

// ErrFailedRequest represents an error when a request fails.
var ErrFailedRequest = errors.New("unexpected status code")

// ErrRateLimitExceeded represents an error when the rate limit is exceeded.
var ErrRateLimitExceeded = errors.New("rate limit exceeded")

// ErrUnsupportedTranslationType represents an error when the translation type is not supported.
var ErrUnsupportedTranslationType = errors.New("unsupported translation type")

// ErrResourceNotFound represents an error when a resource is not found.
var ErrResourceNotFound = errors.New("resource not found")
