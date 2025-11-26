// Package server provides HTTP server utilities with graceful shutdown support.
package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Config holds server configuration.
type Config struct {
	// Port to listen on (default: 8080)
	Port string
	// ShutdownTimeout is the maximum duration to wait for active connections
	// to finish during shutdown (default: 30s)
	ShutdownTimeout time.Duration
	// ReadTimeout is the maximum duration for reading the entire request (default: 30s)
	ReadTimeout time.Duration
	// WriteTimeout is the maximum duration before timing out writes of the response (default: 30s)
	WriteTimeout time.Duration
	// IdleTimeout is the maximum amount of time to wait for the next request (default: 120s)
	IdleTimeout time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig(port string) Config {
	if port == "" {
		port = "8080"
	}
	return Config{
		Port:            port,
		ShutdownTimeout: 30 * time.Second,
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
		IdleTimeout:     120 * time.Second,
	}
}

// Run starts an HTTP server with graceful shutdown support.
// It blocks until SIGTERM or SIGINT is received, then gracefully shuts down.
// The handler is typically a *gin.Engine or any http.Handler.
func Run(handler http.Handler, cfg Config, logger *slog.Logger) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Channel to receive shutdown signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Channel to receive server errors
	serverErr := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Wait for shutdown signal or server error
	select {
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	case sig := <-quit:
		logger.Info("Shutdown signal received", "signal", sig.String())
	}

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	logger.Info("Shutting down server", "timeout", cfg.ShutdownTimeout.String())
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	logger.Info("Server stopped gracefully")
	return nil
}

// RunWithCleanup starts an HTTP server with graceful shutdown and cleanup function support.
// The cleanup function is called after the server has stopped accepting new connections
// but before the function returns. Use this to close database connections, flush buffers, etc.
func RunWithCleanup(handler http.Handler, cfg Config, logger *slog.Logger, cleanup func()) error {
	err := Run(handler, cfg, logger)

	if cleanup != nil {
		logger.Info("Running cleanup")
		cleanup()
		logger.Info("Cleanup completed")
	}

	return err
}
