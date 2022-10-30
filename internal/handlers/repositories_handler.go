package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/internal/model"
	"github.com/internal/server"
	"github.com/internal/services"
	"github.com/internal/shared"
)

type RepositoriesHandler struct {
	service services.IRepositoriesService
}

func NewRepositoriesHandler(service services.IRepositoriesService) *RepositoriesHandler {
	return &RepositoriesHandler{service: service}
}

// GetGroups  godoc
// @Summary     Get all groups
// @Description Needed for create a project in a specific group
// @Tags        Groups
// @Accept      json
// @Produce     json
// @Success     200 {array} model.GroupModel
// @Router      /repositories/groups [get].
func (handler RepositoriesHandler) GetGroups(ctx *fiber.Ctx) error {
	result, err := handler.service.GetGroups()
	if err != nil {
		return err
	}
	return server.SendJSON(ctx, result)
}

func (handler RepositoriesHandler) CreateRepository(ctx *fiber.Ctx) error {
	request := new(model.RepositoryModel)
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

	result, err := handler.service.CreateRepository(request)
	if err != nil {
		return err
	}

	return server.SendJSON(ctx, &model.CreateProjectModel{ID: result})
}

func (handler RepositoriesHandler) GetAppTypes(ctx *fiber.Ctx) error {
	result, err := handler.service.GetAppTypes()
	if err != nil {
		return err
	}

	return server.SendJSON(ctx, result)
}

func (handler RepositoriesHandler) GetApp(ctx *fiber.Ctx) error {
	appName := ctx.Query("app_name")
	err := shared.EnsureNotEmpty(appName, "bad request error, missing app_name")
	if err != nil {
		return err
	}
	result, err := handler.service.GetApp(appName)
	if err != nil {
		return err
	}

	return server.SendJSON(ctx, result)
}
