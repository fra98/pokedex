// Package main is the entry point of the application.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/fra98/pokedex/pkg/flags"
	"github.com/fra98/pokedex/pkg/server"
)

func main() {
	// Initialize options for the application
	opts := flags.Init()

	// Setup the server
	srv := setupServer(opts)

	// Run the server
	if err := runServer(srv, opts); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func setupServer(opts *flags.Options) *http.Server {
	// Setup the Gin engine
	engine := server.SetupEngine()

	// Setup the middlewares
	server.SetupMiddlewares(engine)

	// Register the API endpoints
	server.RegisterEndpoints(engine)

	return &http.Server{
		Addr:         opts.Address,
		Handler:      engine,
		ReadTimeout:  opts.ReadTimeout,
		WriteTimeout: opts.WriteTimeout,
	}
}

func runServer(srv *http.Server, opts *flags.Options) error {
	// Start the server in a separate goroutine to avoid blocking the main thread and handle graceful shutdown
	chanErrors := make(chan error)
	go func() {
		log.Printf("Starting server on %s...", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			chanErrors <- err
		} else {
			chanErrors <- nil
		}
	}()

	// Channel to wait for interrupt signal to gracefully shutdown the server
	chanSignals := make(chan os.Signal, 1)
	signal.Notify(chanSignals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-chanErrors:
		if err != nil {
			return fmt.Errorf("failed to start server: %w", err)
		}
	case sign := <-chanSignals:
		log.Printf("Shutting down server due to signal %q...", sign)

		// Wait for the server to finish processing active requests within the timeout
		ctx, cancel := context.WithTimeout(context.Background(), opts.ShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}

		log.Println("Server gracefully shutdown")
	}

	return nil
}
