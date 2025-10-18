package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func SuccessResponse(c *gin.Context, message string, data any, status ...int) {
	code := http.StatusOK
	if len(status) > 0 {
		code = status[0]
	}

	c.JSON(code, ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, message string, status int, data ...any) {
	var responseData any = nil
	if len(data) > 0 {
		responseData = data[0]
	}

	c.JSON(status, ApiResponse{
		Success: false,
		Message: message,
		Data:    responseData,
	})
}
