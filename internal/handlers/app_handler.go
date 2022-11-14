package handlers

import (
	"net/http"
	"runtime"

	task "github.com/arielsrv/taskpool"
	"github.com/gofiber/fiber/v2"
	"github.com/internal/model"
	"github.com/internal/server"
	"github.com/internal/services"
	"github.com/internal/shared"
)

type AppHandler struct {
	service services.IAppService
	tb      task.Builder
}

func NewAppHandler(service services.IAppService) *AppHandler {
	tb := task.Builder{
		MaxWorkers: runtime.NumCPU() - 1,
	}
	return &AppHandler{
		service: service,
		tb:      tb,
	}
}

// GetGroups  godoc
// @Summary     Get all groups from GitLab
// @Description Needed for create a project in a specific group
// @Tags        Groups
// @Accept      json
// @Produce     json
// @Success     200 {array} model.AppGroupModel
// @Router      /apps/groups [get].
func (h AppHandler) GetGroups(ctx *fiber.Ctx) error {
	result, err := h.service.GetGroups()
	if err != nil {
		return err
	}
	return server.SendJSON(ctx, result)
}

// CreateApp  godoc
// @Summary     Creates an IskayPet Application
// @Tags        Apps
// @Param 		createAppModel    body model.CreateAppModel true "Body params"
// @Router      /apps [post].
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

	err = shared.EnsureInt(request.AppTypeID, "bad request error, invalid app type")
	if err != nil {
		return err
	}

	result, err := h.service.CreateApp(request)
	if err != nil {
		return err
	}

	return server.SendJSON(ctx, result)
}

// GetAppTypes  godoc
// @Summary     Get all application types (backend, frontend, etc.)
// @Description For example: app_name: users-api, key: token, value: hash
// PETS_USERS-API_TOKEN
// @Tags        Apps
// @Router      /apps/types [get].
func (h AppHandler) GetAppTypes(ctx *fiber.Ctx) error {
	result, err := h.service.GetAppTypes()
	if err != nil {
		return err
	}

	return server.SendJSON(ctx, result)
}

// GetApp godoc
// @Summary     Get relevant info for an IskayPet app
// @Tags        Apps
// @Param app_name    query string true "Application name"
// @Router      /apps/search [get].
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

	return server.SendJSON(ctx, result)
}
