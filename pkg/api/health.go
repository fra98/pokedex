package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// IsHealthy returns a 200 OK response if the server is healthy.
func IsHealthy(ctx *gin.Context) {
	res := make(map[string]interface{})
	res["status"] = 200
	res["healthy"] = "OK"
	ctx.JSON(http.StatusOK, res)
}
