package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestLogger middleware logs HTTP requests and responses
func RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Create a unique request ID
		requestID := uuid.New().String()

		// Log format: [REQUEST_ID] METHOD PATH STATUS LATENCY SIZE CLIENT_IP
		return fmt.Sprintf("[%s] %s %s %d %v %s %s\n",
			requestID,
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.BodySize,
			param.ClientIP,
		)
	})
}

// RequestLoggerWithBody middleware logs HTTP requests and responses with request/response body
func RequestLoggerWithBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a unique request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// Start time
		start := time.Now()

		// Read request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// Restore request body for handlers
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create a custom response writer to capture response body
		responseWriter := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = responseWriter

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log request and response
		logRequestResponse(requestID, c, requestBody, responseWriter.body.Bytes(), latency)
	}
}

// responseBodyWriter captures response body
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// logRequestResponse logs the complete request and response
func logRequestResponse(requestID string, c *gin.Context, requestBody, responseBody []byte, latency time.Duration) {
	// Log request
	log.Printf("[%s] REQUEST: %s %s", requestID, c.Request.Method, c.Request.URL.Path)
	log.Printf("[%s] Headers: %v", requestID, c.Request.Header)
	if len(requestBody) > 0 {
		log.Printf("[%s] Request Body: %s", requestID, string(requestBody))
	}

	// Log response
	log.Printf("[%s] RESPONSE: %d", requestID, c.Writer.Status())
	log.Printf("[%s] Response Headers: %v", requestID, c.Writer.Header())
	if len(responseBody) > 0 {
		log.Printf("[%s] Response Body: %s", requestID, string(responseBody))
	}

	// Log timing
	log.Printf("[%s] Latency: %v", requestID, latency)
}

// ErrorLogger middleware logs errors
func ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			requestID, exists := c.Get("request_id")
			if !exists {
				requestID = "unknown"
			}

			for _, err := range c.Errors {
				log.Printf("[%s] ERROR: %s", requestID, err.Error())
			}
		}
	}
}

// PerformanceLogger middleware logs slow requests
func PerformanceLogger(threshold time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		if latency > threshold {
			requestID, exists := c.Get("request_id")
			if !exists {
				requestID = "unknown"
			}

			log.Printf("[%s] SLOW REQUEST: %s %s took %v (threshold: %v)", 
				requestID, c.Request.Method, c.Request.URL.Path, latency, threshold)
		}
	}
}

