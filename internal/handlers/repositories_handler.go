package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/internal/model"
	"github.com/internal/server"
	"github.com/internal/services"
)

type RepositoriesHandler struct {
	service *services.RepositoriesService
}

func NewRepositoriesHandler(service *services.RepositoriesService) *RepositoriesHandler {
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
		return err
	}
	err := handler.service.CreateRepository(request)
	if err != nil {
		return err
	}
	return server.SendString(ctx, "ok")
}
