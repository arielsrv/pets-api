package shared_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/ent/property"

	"github.com/src/shared"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestNewError(t *testing.T) {
	actual := shared.NewError(http.StatusInternalServerError, "nil reference")
	assert.NotNil(t, actual)
	assert.Equal(t, http.StatusInternalServerError, actual.StatusCode)
	assert.Equal(t, "nil reference", actual.Message)
	assert.Equal(t, "nil reference", actual.Error())
}

func TestErrorHandler(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)
	err := shared.ErrorHandler(ctx, errors.New("src server error"))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, ctx.Context().Response.StatusCode())

	var apiError shared.Error
	err = json.Unmarshal(ctx.Response().Body(), &apiError)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, apiError.StatusCode)
	assert.Equal(t, "src server error", apiError.Message)
}

func TestErrorHandler_FiberError(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)
	err := shared.ErrorHandler(ctx, fiber.NewError(http.StatusInternalServerError, "src server error"))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, ctx.Context().Response.StatusCode())

	var apiError shared.Error
	err = json.Unmarshal(ctx.Response().Body(), &apiError)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, apiError.StatusCode)
	assert.Equal(t, "src server error", apiError.Message)
}

func TestErrorHandler_ApiError(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)
	err := shared.ErrorHandler(ctx, shared.NewError(http.StatusInternalServerError, "src server error"))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, ctx.Context().Response.StatusCode())

	var apiError shared.Error
	err = json.Unmarshal(ctx.Response().Body(), &apiError)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, apiError.StatusCode)
	assert.Equal(t, "src server error", apiError.Message)
}

func TestEnsureNotEmpty(t *testing.T) {
	err := shared.EnsureNotEmpty("", "empty value")
	assert.Error(t, err)
	assert.Equal(t, "empty value", err.Error())
}

func TestEnsureInt(t *testing.T) {
	err := shared.EnsureInt(0, "invalid value")
	assert.Error(t, err)
	assert.Equal(t, "invalid value", err.Error())
}

func TestEnsureInt64(t *testing.T) {
	err := shared.EnsureInt64(int64(0), "invalid value")
	assert.Error(t, err)
	assert.Equal(t, "invalid value", err.Error())
}

func TestEnsureNotEmpty_Ok(t *testing.T) {
	err := shared.EnsureNotEmpty("value", "empty value")
	assert.NoError(t, err)
}

func TestEnsureInt_Ok(t *testing.T) {
	err := shared.EnsureInt(1, "invalid value")
	assert.NoError(t, err)
}

func TestEnsureInt64_Ok(t *testing.T) {
	err := shared.EnsureInt64(int64(1), "invalid value")
	assert.NoError(t, err)
}

func TestEnsureEnum(t *testing.T) {
	err := shared.EnsureEnum(property.Backend, property.AppTypeValues, "invalid value")
	assert.NoError(t, err)
}

func TestEnsureEnum_Err(t *testing.T) {
	err := shared.EnsureEnum(0, property.AppTypeValues, "invalid value")
	assert.Error(t, err)
	assert.Equal(t, "invalid value", err.Error())
}
