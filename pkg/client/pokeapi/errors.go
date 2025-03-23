package pokeapi

import "errors"

// ErrFailedRequest represents an error when a request fails.
var ErrFailedRequest = errors.New("unexpected status code")
