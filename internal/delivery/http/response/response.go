package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// StandardResponse standard API response structure
type StandardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Error   *ErrorData  `json:"error,omitempty"`
}

// Response is an alias for StandardResponse for Swagger documentation
type Response = StandardResponse

// ErrorData detailed error structure
type ErrorData struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Meta standard pagination metadata
type Meta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// Success sends a standard success response
func Success(c *gin.Context, statusCode int, message string, data interface{}, meta interface{}) {
	c.JSON(statusCode, StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Error sends a standard error response
func Error(c *gin.Context, statusCode int, errorCode string, message string, details interface{}) {
	c.JSON(statusCode, StandardResponse{
		Success: false,
		Error: &ErrorData{
			Code:    errorCode,
			Message: message,
			Details: details,
		},
	})
}

// BadRequest helper for 400 Bad Request
func BadRequest(c *gin.Context, message string, details interface{}) {
	Error(c, http.StatusBadRequest, "BAD_REQUEST", message, details)
}

// Unauthorized helper for 401 Unauthorized
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

// Forbidden helper for 403 Forbidden
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, "FORBIDDEN", message, nil)
}

// NotFound helper for 404 Not Found
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, "NOT_FOUND", message, nil)
}

// InternalServerError helper for 500 Internal Server Error
func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message, nil)
}
