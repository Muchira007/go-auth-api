package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func LogRequest(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Log after request
		latency := time.Since(start)
		fmt.Printf("Request took %v\n", latency)
		fmt.Printf("Request path: %s\n", c.Request.URL.Path)
}
