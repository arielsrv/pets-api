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
	err := handler.service.CreateRepository(request)
	if err != nil {
		return err
	}
	return server.SendString(ctx, "ok")
}
