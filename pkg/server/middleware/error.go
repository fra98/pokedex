package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fra98/pokedex/pkg/server/httperror"
)

// ErrorHandler is a middleware that handles errors and returns a JSON response.
// If the error is of type errors.HTTP, it will return the error message and status code.
// If the error is of any other type, it will return a generic error message and status code 500.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, err := range c.Errors {
			var httperr httperror.HTTPError
			switch {
			case errors.As(err.Err, &httperr):
				c.AbortWithStatusJSON(httperr.StatusCode, httperr)
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"message": "Internal Server Error"})
			}
		}
	}
}
