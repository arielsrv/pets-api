package server

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Error struct {
	StatusCode int    `json:"status_code,omitempty"`
	Message    string `json:"message,omitempty"`
}

func NewError(statusCode int, message string) *Error {
	return &Error{StatusCode: statusCode, Message: message}
}

func (e *Error) Error() string {
	return e.Message
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	var e = new(Error)

	var fiberError *fiber.Error
	var apiError *Error

	switch {
	case errors.As(err, &fiberError):
		{
			e.StatusCode = fiberError.Code
			e.Message = fiberError.Message
		}
	case errors.As(err, &apiError):
		{
			e.StatusCode = apiError.StatusCode
			e.Message = apiError.Message
		}
	default:
		e.StatusCode = http.StatusInternalServerError
		e.Message = err.Error()
	}

	ctx.Status(e.StatusCode)
	return ctx.JSON(e)
}

func EnsureNotEmpty(value string, message string) error {
	if value == "" {
		return NewError(http.StatusBadRequest, message)
	}
	return nil
}

func EnsureInt(value int, message string) error {
	if value < 1 {
		return NewError(http.StatusBadRequest, message)
	}
	return nil
}

func EnsureInt64(value int64, message string) error {
	if value < 1 {
		return NewError(http.StatusBadRequest, message)
	}
	return nil
}

func EnsureEnum[T comparable](value T, elements []T, message string) error {
	valid := false
	for _, element := range elements {
		if value == element {
			valid = true
			continue
		}
	}
	if !valid {
		return NewError(http.StatusBadRequest, message)
	}
	return nil
}