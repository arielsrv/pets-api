package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/internal/server"
	"github.com/internal/services"
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
// @Summary     Check if the instance is healthy or unhealthy
// @Description Health
// @Tags        Check
// @Success     200
// @Produce     plain
// @Success     200 {string} string "pong"
// @Router      /ping [get].
func (h PingHandler) Ping(ctx *fiber.Ctx) error {
	result := h.pingService.Ping()
	return server.SendString(ctx, result)
}
