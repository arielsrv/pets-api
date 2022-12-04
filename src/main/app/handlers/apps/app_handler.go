package apps

import (
	"net/http"

	"github.com/src/main/app/model"

	"github.com/src/main/app/helpers/ensure"

	"github.com/src/main/app/services/apps"

	"github.com/gofiber/fiber/v2"
	"github.com/src/main/app/server"
)

type AppHandler struct {
	service apps.IAppService
}

func NewAppHandler(service apps.IAppService) *AppHandler {
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
	request := new(model.CreateAppRequest)
	if err := ctx.BodyParser(request); err != nil {
		return server.NewError(http.StatusBadRequest, "bad request error")
	}

	err := ensure.NotEmpty(request.Name, "bad request error, missing name")
	if err != nil {
		return err
	}

	err = ensure.Int64(request.GroupID, "bad request error, invalid group id")
	if err != nil {
		return err
	}

	err = ensure.Enum(request.AppTypeID, "bad request error, invalid app type")
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
	err := ensure.NotEmpty(appName, "bad request error, missing app_name")
	if err != nil {
		return err
	}

	result, err := h.service.GetAppByName(appName)
	if err != nil {
		return err
	}

	return server.SendOk(ctx, result)
}
