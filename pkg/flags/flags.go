// Package flags provides flags to configure the server.
package flags

import (
	"time"

	"github.com/spf13/pflag"
)

// Init initializes flags to configure the server.
func Init() *Options {
	opts := NewOptions()

	pflag.StringVar(&opts.Address, "address", ":8080", "Address to listen on")
	pflag.DurationVar(&opts.ReadTimeout, "read-timeout", 10*time.Second, "Read timeout for the server")
	pflag.DurationVar(&opts.WriteTimeout, "write-timeout", 10*time.Second, "Write timeout for the server")
	pflag.DurationVar(&opts.ShutdownTimeout, "shutdown-timeout", 10*time.Second, "Graceful shutdown timeout for the server")
	pflag.BoolVar(&opts.DisableCache, "disable-cache", false, "Disable caching")
	pflag.DurationVar(&opts.CacheTimeoutExpiration, "cache-timeout-expiration", 1*time.Hour, "Cache timeout expiration")
	pflag.DurationVar(&opts.CacheCleanupInterval, "cache-cleanup-interval", 24*time.Hour, "Cache cleanup interval")

	pflag.Parse()

	return opts
}
