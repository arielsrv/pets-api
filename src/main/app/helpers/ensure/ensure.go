package ensure

import (
	"net/http"

	"github.com/arielsrv/pets-api/src/main/app/server"
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

type SafeEnum interface {
	IsValid() bool
}

func Enum(safeEnum SafeEnum, message string) error {
	if !safeEnum.IsValid() {
		return server.NewError(http.StatusBadRequest, message)
	}

	return nil
}
