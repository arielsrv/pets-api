package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/src/main/app/server"
	"github.com/src/main/app/services"
)

type SnippetHandler struct {
	snippetService services.ISnippetService
}

func NewSnippetHandler(snippetService services.ISnippetService) *SnippetHandler {
	return &SnippetHandler{snippetService: snippetService}
}

func (h SnippetHandler) GetSnippet(ctx *fiber.Ctx) error {
	secretID, err := ctx.ParamsInt("secretID")
	if err != nil {
		return err
	}

	err = server.EnsureInt64(int64(secretID), "bad request error, invalid secret id")
	if err != nil {
		return err
	}

	model, err := h.snippetService.GetSecrets(int64(secretID))
	if err != nil {
		return err
	}

	return ctx.Render("snippets/index", fiber.Map{
		"Model": model,
	})
}
