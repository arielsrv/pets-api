package server_test

import (
	"encoding/json"
	"errors"
	"github.com/src/main/app/server"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/src/main/app/ent/property"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestNewError(t *testing.T) {
	actual := server.NewError(http.StatusInternalServerError, "nil reference")
	assert.NotNil(t, actual)
	assert.Equal(t, http.StatusInternalServerError, actual.StatusCode)
	assert.Equal(t, "nil reference", actual.Message)
	assert.Equal(t, "nil reference", actual.Error())
}

func TestErrorHandler(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)
	err := server.ErrorHandler(ctx, errors.New("src server error"))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, ctx.Context().Response.StatusCode())

	var apiError server.Error
	err = json.Unmarshal(ctx.Response().Body(), &apiError)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, apiError.StatusCode)
	assert.Equal(t, "src server error", apiError.Message)
}

func TestErrorHandler_FiberError(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)
	err := server.ErrorHandler(ctx, fiber.NewError(http.StatusInternalServerError, "src server error"))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, ctx.Context().Response.StatusCode())

	var apiError server.Error
	err = json.Unmarshal(ctx.Response().Body(), &apiError)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, apiError.StatusCode)
	assert.Equal(t, "src server error", apiError.Message)
}

func TestErrorHandler_ApiError(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)
	err := server.ErrorHandler(ctx, server.NewError(http.StatusInternalServerError, "src server error"))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, ctx.Context().Response.StatusCode())

	var apiError server.Error
	err = json.Unmarshal(ctx.Response().Body(), &apiError)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, apiError.StatusCode)
	assert.Equal(t, "src server error", apiError.Message)
}

func TestEnsureNotEmpty(t *testing.T) {
	err := server.EnsureNotEmpty("", "empty value")
	assert.Error(t, err)
	assert.Equal(t, "empty value", err.Error())
}

func TestEnsureInt(t *testing.T) {
	err := server.EnsureInt(0, "invalid value")
	assert.Error(t, err)
	assert.Equal(t, "invalid value", err.Error())
}

func TestEnsureInt64(t *testing.T) {
	err := server.EnsureInt64(int64(0), "invalid value")
	assert.Error(t, err)
	assert.Equal(t, "invalid value", err.Error())
}

func TestEnsureNotEmpty_Ok(t *testing.T) {
	err := server.EnsureNotEmpty("value", "empty value")
	assert.NoError(t, err)
}

func TestEnsureInt_Ok(t *testing.T) {
	err := server.EnsureInt(1, "invalid value")
	assert.NoError(t, err)
}

func TestEnsureInt64_Ok(t *testing.T) {
	err := server.EnsureInt64(int64(1), "invalid value")
	assert.NoError(t, err)
}

func TestEnsureEnum(t *testing.T) {
	err := server.EnsureEnum(property.Backend, property.AppTypeValues, "invalid value")
	assert.NoError(t, err)
}

func TestEnsureEnum_Err(t *testing.T) {
	err := server.EnsureEnum(0, property.AppTypeValues, "invalid value")
	assert.Error(t, err)
	assert.Equal(t, "invalid value", err.Error())
}
