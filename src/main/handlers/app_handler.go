package handlers

import (
	"net/http"

	"github.com/src/main/model"
	"github.com/src/main/server"
	"github.com/src/main/services"
	"github.com/src/main/shared"

	"github.com/ent/property"

	"github.com/gofiber/fiber/v2"
)

type AppHandler struct {
	service services.IAppService
}

func NewAppHandler(service services.IAppService) *AppHandler {
	return &AppHandler{
		service: service,
	}
}

func (h AppHandler) GetGroups(ctx *fiber.Ctx) error {
	result, err := h.service.GetGroups()
	if err != nil {
		return err
	}
	return server.SendOk(ctx, result)
}

func (h AppHandler) CreateApp(ctx *fiber.Ctx) error {
	request := new(model.CreateAppModel)
	if err := ctx.BodyParser(request); err != nil {
		return shared.NewError(http.StatusBadRequest, "bad request error")
	}

	err := shared.EnsureNotEmpty(request.Name, "bad request error, missing name")
	if err != nil {
		return err
	}

	err = shared.EnsureInt64(request.GroupID, "bad request error, invalid group id")
	if err != nil {
		return err
	}

	err = shared.EnsureEnum(request.AppTypeID, property.AppTypeValues, "bad request error, invalid app type")
	if err != nil {
		return err
	}

	result, err := h.service.CreateApp(request)
	if err != nil {
		return err
	}

	return server.SendCreated(ctx, result)
}

func (h AppHandler) GetAppTypes(ctx *fiber.Ctx) error {
	result, err := h.service.GetAppTypes()
	if err != nil {
		return err
	}

	return server.SendOk(ctx, result)
}

func (h AppHandler) GetApp(ctx *fiber.Ctx) error {
	appName := ctx.Query("app_name")
	err := shared.EnsureNotEmpty(appName, "bad request error, missing app_name")
	if err != nil {
		return err
	}
	result, err := h.service.GetAppByName(appName)
	if err != nil {
		return err
	}

	return server.SendOk(ctx, result)
}
