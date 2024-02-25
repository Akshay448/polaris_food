package middlewares

import (
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extract the token from the Authorization header
		// 2. Parse the token
		// 3. Validate the token
		// 4. If valid, proceed to the next handler. Otherwise, return an error.
	}
}
