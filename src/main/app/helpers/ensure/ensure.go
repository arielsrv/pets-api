package ensure

import (
	"net/http"

	"github.com/src/main/app/server"
)

func NotEmpty(value string, message string) error {
	if value == "" {
		return server.NewError(http.StatusBadRequest, message)
	}
	return nil
}

func Int(value int, message string) error {
	if value < 1 {
		return server.NewError(http.StatusBadRequest, message)
	}
	return nil
}

func Int64(value int64, message string) error {
	if value < 1 {
		return server.NewError(http.StatusBadRequest, message)
	}
	return nil
}

func Enum[T comparable](value T, elements []T, message string) error {
	valid := false
	for _, element := range elements {
		if value == element {
			valid = true
			continue
		}
	}
	if !valid {
		return server.NewError(http.StatusBadRequest, message)
	}
	return nil
}
