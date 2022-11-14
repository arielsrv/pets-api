package handlers

import (
	"github.com/internal/model"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/internal/server"
	"github.com/internal/services"
	"github.com/internal/shared"
)

type SecretHandler struct {
	appService    services.IAppService
	secretService services.ISecretService
}

func NewSecretHandler(appService services.IAppService, secretService services.ISecretService) *SecretHandler {
	return &SecretHandler{
		appService:    appService,
		secretService: secretService,
	}
}

func (h SecretHandler) CreateSecret(ctx *fiber.Ctx) error {
	appId, err := ctx.ParamsInt("appId")
	if err != nil {
		return err
	}

	appModel, err := h.appService.GetAppById(int64(appId))
	if err != nil {
		return err
	}

	request := new(model.CreateSecretRequestModel)
	if err = ctx.BodyParser(request); err != nil {
		return shared.NewError(http.StatusBadRequest, "bad request error, missing key and value properties")
	}

	err = shared.EnsureNotEmpty(request.Key, "bad request error, missing key")
	if err != nil {
		return err
	}

	err = shared.EnsureNotEmpty(request.Value, "bad request error, missing value")
	if err != nil {
		return err
	}

	result, err := h.secretService.SaveSecret(appModel.ID, request)
	if err != nil {
		return err
	}

	result.RelativeUrl = ctx.BaseURL() + result.RelativeUrl
	ctx.Status(201)

	return server.SendJSON(ctx, result)
}
