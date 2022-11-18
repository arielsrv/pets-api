package handlers

import (
	"net/http"

	"github.com/src/main/services"

	"github.com/src/main/model"
	"github.com/src/main/server"
	"github.com/src/main/shared"

	"github.com/gofiber/fiber/v2"
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

	request := new(model.CreateAppSecretModel)
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

	result, err := h.secretService.SaveSecret(int64(appId), request)
	if err != nil {
		return err
	}

	return server.SendCreated(ctx, result)
}
