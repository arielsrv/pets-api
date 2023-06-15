package secrets

import (
	"net/http"

	"github.com/arielsrv/pets-api/src/main/app/model"

	"github.com/arielsrv/pets-api/src/main/app/helpers/ensure"

	"github.com/arielsrv/pets-api/src/main/app/services/apps"
	"github.com/arielsrv/pets-api/src/main/app/services/secrets"

	"github.com/arielsrv/pets-api/src/main/app/server"
	"github.com/gofiber/fiber/v2"
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
// @Param		createAppSecretModel	body model.CreateSecretRequest  true "Body params"
// @Param		appID  path int true "Pet ID"
// @Success		200  {array}   model.CreateSecretResponse
// @Failure		404  {object}  server.Error "App not found"
// @Failure		409  {object}  server.Error "Key already exist"
// @Failure		500  {object}  server.Error "Internal server error"
// @Router		/apps/{appID}/secrets [post].
func (h SecretHandler) CreateSecret(ctx *fiber.Ctx) error {
	appID, err := ctx.ParamsInt("appID")
	if err != nil {
		return err
	}

	request := new(model.CreateSecretRequest)
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
