package utils

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrNotFound       = errors.New("resource not found")
	ErrBadRequest     = errors.New("invalid request")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrInternalServer = errors.New("internal server error")
	ErrDuplicate      = errors.New("resource already exists")
	ErrValidation     = errors.New("validation error")
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"
	details := err.Error()

	// Fiber error handling
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// Custom error handling
	switch err {
	case ErrNotFound:
		code = fiber.StatusNotFound
		message = "Resource Not Found"
	case ErrBadRequest:
		code = fiber.StatusBadRequest
		message = "Bad Request"
	case ErrUnauthorized:
		code = fiber.StatusUnauthorized
		message = "Unauthorized"
	case ErrForbidden:
		code = fiber.StatusForbidden
		message = "Forbidden"
	case ErrDuplicate:
		code = fiber.StatusConflict
		message = "Resource Already Exists"
	case ErrValidation:
		code = fiber.StatusBadRequest
		message = "Validation Error"
	}

	return ErrorResponse(c, code, message, details)
}

type ValidationErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
