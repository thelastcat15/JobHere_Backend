package utils

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func SuccessResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Status:  statusCode,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *fiber.Ctx, statusCode int, message string, err interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Status:  statusCode,
		Message: message,
		Error:   err,
	})
}

type PaginatedData struct {
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int64       `json:"total_pages"`
	Items      interface{} `json:"items"`
}

func PaginatedResponse(c *fiber.Ctx, statusCode int, message string, data PaginatedData) error {
	return c.Status(statusCode).JSON(Response{
		Status:  statusCode,
		Message: message,
		Data:    data,
	})
}
