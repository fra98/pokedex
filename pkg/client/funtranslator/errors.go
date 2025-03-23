package funtranslator

import "errors"

// ErrFailedRequest represents an error when a request fails.
var ErrFailedRequest = errors.New("unexpected status code")

// ErrRateLimitExceeded represents an error when the rate limit is exceeded.
var ErrRateLimitExceeded = errors.New("rate limit exceeded")

// ErrUnsupportedTranslationType represents an error when the translation type is not supported.
var ErrUnsupportedTranslationType = errors.New("unsupported translation type")
