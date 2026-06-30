package response

import (
	"task-api/appError"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success    bool        `json:"success"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success:    true,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}

func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Success:    false,
		StatusCode: statusCode,
		Message:    message,
		Data:       nil,
	})
}

func ValidationError(c *gin.Context, statusCode int, message string, errors interface{}) {
	c.JSON(statusCode, Response{
		Success:    false,
		StatusCode: statusCode,
		Message:    message,
		Data:       errors,
	})
}

func FromAppError(c *gin.Context, appErr *appError.AppError) {
	c.JSON(appErr.Status, Response{
		Success:    false,
		StatusCode: appErr.Status,
		Message:    appErr.Message,
		Data:       nil,
	})
}
