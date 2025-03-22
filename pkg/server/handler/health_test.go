package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/fra98/pokedex/pkg/server/handler"
)

func TestIsHealthySuccess(t *testing.T) {
	t.Parallel()

	// Create a httptest recorder
	responseRecorder := httptest.NewRecorder()

	// Create a gin text context for the above recorder to get context and gin engine
	ctx, engine := gin.CreateTestContext(responseRecorder)

	// Register the healthcheck endpoint to the gin engine
	engine.GET("/health", handler.IsHealthy)

	// Create a test request for the above registered endpoint
	ctx.Request = httptest.NewRequest(http.MethodGet, "/health", http.NoBody)

	// Below line shows how to set headers if necessary
	ctx.Request.Header.Set("X-Agentname", "agent name")

	// Test the endpoint
	engine.ServeHTTP(responseRecorder, ctx.Request)

	// Assert if response http status is as expected
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}
