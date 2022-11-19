package handlers

import (
	"net/http"

	"github.com/src/main/app/services"

	"github.com/gofiber/fiber/v2"
	"github.com/src/main/app/model"
	"github.com/src/main/app/server"
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
		return server.NewError(http.StatusBadRequest, "bad request error, missing key and value properties")
	}

	err = server.EnsureNotEmpty(request.Key, "bad request error, missing key")
	if err != nil {
		return err
	}

	err = server.EnsureNotEmpty(request.Value, "bad request error, missing value")
	if err != nil {
		return err
	}

	result, err := h.secretService.SaveSecret(int64(appId), request)
	if err != nil {
		return err
	}

	return server.SendCreated(ctx, result)
}