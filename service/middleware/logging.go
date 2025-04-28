package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		logMessage := fmt.Sprintf(
			"Method: %s, Path: %s, Status: %d, Duration: %s, ClientIP: %s, UserAgent: %s",
			c.Request.Method,
			c.FullPath(),
			c.Writer.Status(),
			duration,
			c.ClientIP(),
			c.Request.UserAgent(),
		)

		log.Println(logMessage)
	}
}
