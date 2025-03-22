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

	pflag.Parse()

	return opts
}
