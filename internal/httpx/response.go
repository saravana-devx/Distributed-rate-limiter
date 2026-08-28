// Package httpx provides a consistent JSON response envelope for all HTTP
// handlers. Every response carries a "success" boolean and a human-readable
// "message"; successful responses may additionally include a "data" payload.
package httpx

import "github.com/gin-gonic/gin"

// Success writes a success envelope: {"success": true, "message": ..., "data": ...}.
// When data is nil the "data" field is omitted (e.g. delete/logout responses).
func Success(c *gin.Context, status int, message string, data any) {
	body := gin.H{
		"success": true,
		"message": message,
	}
	if data != nil {
		body["data"] = data
	}
	c.JSON(status, body)
}

// Error writes a failure envelope: {"success": false, "message": ...}.
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"message": message,
	})
}

// AbortError writes a failure envelope and aborts the handler chain. Use this
// from middleware where downstream handlers must not run.
func AbortError(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, gin.H{
		"success": false,
		"message": message,
	})
}
