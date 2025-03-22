package server

import (
	"github.com/gin-gonic/gin"

	"github.com/fra98/pokedex/pkg/server/handler"
	"github.com/fra98/pokedex/pkg/server/middleware"
)

// SetupEngine sets up and returns a new Gin engine.
func SetupEngine() *gin.Engine {
	return gin.Default() // default Gin engine with Logger and Recovery middleware already attached
}

// SetupMiddlewares sets up the middlewares for the server engine.
func SetupMiddlewares(r *gin.Engine) {
	// Register the error handler middleware
	r.Use(middleware.ErrorHandler())
}

// RegisterEndpoints registers the endpoints of the API to the server engine.
func RegisterEndpoints(r *gin.Engine) {
	// Group the endpoints by version
	v1 := r.Group("/v1")

	// Health check
	v1.GET("/health", handler.IsHealthy)
}
