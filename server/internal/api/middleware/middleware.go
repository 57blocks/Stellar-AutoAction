package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func Header() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func Response() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Next()
	}
}

func Error() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check if any error occurred
		if len(c.Errors) > 0 {
			c.Header("Content-Type", "application/json")
			// Log or handle the errors here as needed
			// For example, return a JSON response with the error
			c.JSON(
				// TODO: update the status code of error
				http.StatusInternalServerError,
				gin.H{
					"status":  http.StatusInternalServerError,
					"message": c.Errors.Last().Error(),
					"notes":   "internal error",
				},
			)

			// Prevent calling any subsequent handlers
			c.Abort()
		}
	}
}
