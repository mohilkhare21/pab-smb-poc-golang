package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS middleware handles Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// CORSCustom middleware allows custom CORS configuration
func CORSCustom(origins []string, methods []string, headers []string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Set allowed origins
		origin := c.Request.Header.Get("Origin")
		if len(origins) == 0 {
			c.Header("Access-Control-Allow-Origin", "*")
		} else {
			for _, allowedOrigin := range origins {
				if origin == allowedOrigin {
					c.Header("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		// Set allowed credentials
		c.Header("Access-Control-Allow-Credentials", "true")

		// Set allowed headers
		if len(headers) == 0 {
			c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		} else {
			headerStr := ""
			for i, header := range headers {
				if i > 0 {
					headerStr += ", "
				}
				headerStr += header
			}
			c.Header("Access-Control-Allow-Headers", headerStr)
		}

		// Set allowed methods
		if len(methods) == 0 {
			c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		} else {
			methodStr := ""
			for i, method := range methods {
				if i > 0 {
					methodStr += ", "
				}
				methodStr += method
			}
			c.Header("Access-Control-Allow-Methods", methodStr)
		}

		// Set max age for preflight requests
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}
