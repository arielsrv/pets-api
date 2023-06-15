package handlers

import (
	"github.com/arielsrv/pets-api/src/main/app/server"
	"github.com/arielsrv/pets-api/src/main/app/services"
	"github.com/gofiber/fiber/v2"
)

type IPingHandler interface {
	Ping(ctx *fiber.Ctx) error
}

type PingHandler struct {
	pingService services.IPingService
}

func NewPingHandler(pingService services.IPingService) *PingHandler {
	return &PingHandler{
		pingService: pingService,
	}
}

// Ping godoc
//
//	@Summary		Check if the instance is healthy or unhealthy
//	@Description	Health
//	@Tags			health
//	@Success		200
//	@Produce		plain
//	@Success		200	{string}	string	"pong"
//	@Router			/ping [get].
func (h PingHandler) Ping(ctx *fiber.Ctx) error {
	result := h.pingService.Ping()
	return server.SendString(ctx, result)
}
