package handler

import (
	"github.com/RevinB/mira_server/config"
	"github.com/RevinB/mira_server/data"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	data   data.Store
	config config.Config
}

func ImplHandler(r *fiber.App) {
	{
		//userGroup := r.Group("/user")

	}
}

func (h *Handler) Data() data.Store {
	return h.data
}

func (h *Handler) Config() config.Config {
	return h.config
}
