package handlers

import (
	"fmt"
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

// CreateSecret  godoc
// @Summary     Creates secret for application
// @Tags        Apps
// @Param appId    path int true "App ID"
// @Param 		createAppSecretModel    body model.CreateAppSecretModel true "Body params"
// @Router      /apps/{appId}/secrets [post].
func (h AppHandler) CreateSecret(ctx *fiber.Ctx) error {
	appId, err := ctx.ParamsInt("appId")
	if err != nil {
		return err
	}

	appModel, err := h.service.GetAppById(int64(appId))
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

	// result, err = secrets.CreateSecret(appModel.Name, request.Key, request.Value)
	// if err != nil {
	//     return err
	// }
	// service.PutSecret( ... result ...)

	// ej: appName: customers-api, key: mirakl_token, value: abcdef12345
	// secrets.GetSecret(PETS_CUSTOMERS-API_MIRAKL_TOKEN) // PREFIX_APPNAME-KEY
	// result

	result := new(model.AppSecretModel)
	result.Key = request.Key //nolint:nolintlint,govet

	var snippets []model.SnippetModel
	var snippet1 model.SnippetModel

	// @todo: download snippet code from any repo so parse and, replace variables and servers url
	codeUrl := fmt.Sprintf("%s/apps/%d/secrets/snippets", ctx.BaseURL(), appModel.ID)
	// apps/123/secrets
	// secrets/1/secrets/1
	snippet1.CodeUrl = codeUrl

	snippets = append(snippets, snippet1)

	result.Snippets = snippets

	ctx.Status(201)
	return server.SendJSON(ctx, result)
}
