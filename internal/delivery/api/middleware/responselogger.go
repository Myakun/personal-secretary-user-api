package middleware

import (
	"bytes"
	"fmt"

	"github.com/gin-gonic/gin"
)

// CustomResponseWriter captures the response body
type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)                  // capture body
	return w.ResponseWriter.Write(b) // write to actual response
}

func ResponseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Replace writer with custom one
		blw := &CustomResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Proceed with request
		c.Next()

		// Log status and response body
		statusCode := c.Writer.Status()
		responseBody := blw.body.String()

		fmt.Println("RESPONSE LOGGER: STATUS: %d RESPONSE: %s", statusCode, responseBody)
	}
}
