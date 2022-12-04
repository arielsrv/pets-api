package secrets

import (
	"net/http"

	"github.com/src/main/app/model"

	"github.com/src/main/app/helpers/ensure"

	"github.com/src/main/app/services/apps"
	"github.com/src/main/app/services/secrets"

	"github.com/gofiber/fiber/v2"
	"github.com/src/main/app/server"
)

type SecretHandler struct {
	appService    apps.IAppService
	secretService secrets.ISecretService
}

func NewSecretHandler(appService apps.IAppService, secretService secrets.ISecretService) *SecretHandler {
	return &SecretHandler{
		appService:    appService,
		secretService: secretService,
	}
}

// CreateSecret godoc
// @Summary		Creates the secret
// @Description	Get snippet key, conflict if secret already exist.
// @Tags		secrets
// @Accept		json
// @Produce		json
// @Param		createAppSecretModel	body model.CreateAppSecretModel  true "Body params"
// @Param		appID  path int true "Pet ID"
// @Success		200  {array}   model.AppSecretModel
// @Failure		404  {object}  server.Error "App not found"
// @Failure		409  {object}  server.Error "Key already exist"
// @Failure		500  {object}  server.Error "Internal server error"
// @Router		/apps/{appID}/secrets [post].
func (h SecretHandler) CreateSecret(ctx *fiber.Ctx) error {
	appID, err := ctx.ParamsInt("appID")
	if err != nil {
		return err
	}

	request := new(model.CreateAppSecretModel)
	if err = ctx.BodyParser(request); err != nil {
		return server.NewError(http.StatusBadRequest, "bad request error, missing key and value properties")
	}

	err = ensure.NotEmpty(request.Key, "bad request error, missing key")
	if err != nil {
		return err
	}

	err = ensure.NotEmpty(request.Value, "bad request error, missing value")
	if err != nil {
		return err
	}

	result, err := h.secretService.CreateSecret(int64(appID), request)
	if err != nil {
		return err
	}

	return server.SendCreated(ctx, result)
}
