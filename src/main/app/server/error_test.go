package server_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/arielsrv/pets-api/src/main/app/ent/property"
	"github.com/arielsrv/pets-api/src/main/app/helpers/ensure"
	"github.com/arielsrv/pets-api/src/main/app/server"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	err := server.ErrorHandler(ctx, errors.New("src ensure error"))
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, ctx.Context().Response.StatusCode())

	var apiError server.Error
	err = json.Unmarshal(ctx.Response().Body(), &apiError)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, apiError.StatusCode)
	assert.Equal(t, "src ensure error", apiError.Message)
}

func TestErrorHandler_FiberError(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)
	err := server.ErrorHandler(ctx, fiber.NewError(http.StatusInternalServerError, "src ensure error"))
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, ctx.Context().Response.StatusCode())

	var apiError server.Error
	err = json.Unmarshal(ctx.Response().Body(), &apiError)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, apiError.StatusCode)
	assert.Equal(t, "src ensure error", apiError.Message)
}

func TestErrorHandler_ApiError(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)
	err := server.ErrorHandler(ctx, server.NewError(http.StatusInternalServerError, "src ensure error"))
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, ctx.Context().Response.StatusCode())

	var apiError server.Error
	err = json.Unmarshal(ctx.Response().Body(), &apiError)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, apiError.StatusCode)
	assert.Equal(t, "src ensure error", apiError.Message)
}

func TestEnsureNotEmpty(t *testing.T) {
	err := ensure.NotEmpty("", "empty value")
	require.Error(t, err)
	assert.Equal(t, "empty value", err.Error())
}

func TestEnsureInt(t *testing.T) {
	err := ensure.Int(0, "invalid value")
	require.Error(t, err)
	assert.Equal(t, "invalid value", err.Error())
}

func TestEnsureInt64(t *testing.T) {
	err := ensure.Int64(int64(0), "invalid value")
	require.Error(t, err)
	assert.Equal(t, "invalid value", err.Error())
}

func TestEnsureNotEmpty_Ok(t *testing.T) {
	err := ensure.NotEmpty("value", "empty value")
	require.NoError(t, err)
}

func TestEnsureInt_Ok(t *testing.T) {
	err := ensure.Int(1, "invalid value")
	require.NoError(t, err)
}

func TestEnsureInt64_Ok(t *testing.T) {
	err := ensure.Int64(int64(1), "invalid value")
	require.NoError(t, err)
}

func TestEnsureEnum(t *testing.T) {
	err := ensure.Enum(property.Backend, "invalid value")
	require.NoError(t, err)
}
