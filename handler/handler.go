package handler

import (
	"github.com/RevinB/mira_server/config"
	"github.com/RevinB/mira_server/data"
	"github.com/RevinB/mira_server/router"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	data   data.Store
	config config.Config
}

func NewHandler(data data.Store, cfg config.Config) Handler {
	return Handler{
		data:   data,
		config: cfg,
	}
}

func (h *Handler) ImplHandler(r *fiber.App) {
	{
		userGroup := r.Group("/user")
		userGroup.Post("/reset", h.UserResetSecret)

		userGroup.Use(router.GetJwtHandler())
	}
}

func (h *Handler) Data() data.Store {
	return h.data
}

func (h *Handler) Config() config.Config {
	return h.config
}
