package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/src/main/app/services"
	"github.com/src/main/app/shared"
)

type SnippetHandler struct {
	snippetService services.ISnippetService
}

func NewSnippetHandler(snippetService services.ISnippetService) *SnippetHandler {
	return &SnippetHandler{snippetService: snippetService}
}

func (h SnippetHandler) GetSnippet(ctx *fiber.Ctx) error {
	secretId, err := ctx.ParamsInt("secretId")
	if err != nil {
		return err
	}

	err = shared.EnsureInt64(int64(secretId), "bad request error, invalid secret id")
	if err != nil {
		return err
	}

	model, err := h.snippetService.GetSecrets(int64(secretId))
	if err != nil {
		return err
	}

	return ctx.Render("snippets/index", fiber.Map{
		"Model": model,
	})
}
