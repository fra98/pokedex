package flags

import (
	"time"
)

// NewOptions returns a new Options instance.
func NewOptions() *Options {
	return &Options{}
}

// Options contains the server options.
type Options struct {
	// Server options
	Address         string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	// Cache options
	DisableCache           bool
	CacheTimeoutExpiration time.Duration
	CacheCleanupInterval   time.Duration
}
