package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func StripTrailingSlash() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path != "/" {
			c.Request.URL.Path = stripTrailingSlash(c.Request.URL.Path)
		}
		c.Next()
	}
}

func stripTrailingSlash(value string) string {
	if len(value) > 1 && value[len(value)-1] == '/' {
		return value[:len(value)-1]
	}
	return value
}

// other suggested middleware:
// - CORS
// - Rate limiting
// - Logging
// - Error handling

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin, Accept")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				log.Printf("Error: %s", e.Error())
			}
		}
	}
}
