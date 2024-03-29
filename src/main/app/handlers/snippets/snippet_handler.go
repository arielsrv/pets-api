package snippets

import (
	"github.com/arielsrv/pets-api/src/main/app/helpers/ensure"
	"github.com/arielsrv/pets-api/src/main/app/services/snippets"
	"github.com/gofiber/fiber/v2"
)

type SnippetHandler struct {
	snippetService snippets.ISnippetService
}

func NewSnippetHandler(snippetService snippets.ISnippetService) *SnippetHandler {
	return &SnippetHandler{snippetService: snippetService}
}

func (h SnippetHandler) GetSnippet(ctx *fiber.Ctx) error {
	secretID, err := ctx.ParamsInt("secretID")
	if err != nil {
		return err
	}

	err = ensure.Int64(int64(secretID), "bad request error, invalid secret id")
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
