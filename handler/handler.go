package handler

import (
	"github.com/RevinB/mira_server/config"
	"github.com/RevinB/mira_server/data"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	data   data.Store
	config config.Config
	aws    *session.Session
}

func NewHandler(data data.Store, cfg config.Config, a *session.Session) Handler {
	return Handler{
		data:   data,
		config: cfg,
		aws:    a,
	}
}

func (h *Handler) ImplHandler(r *fiber.App) {
	{
		userGroup := r.Group("/user")
		userGroup.Post("/reset/:id", h.UserResetSecret)

		userGroup.Use(h.JwtMiddleware)

		userGroup.Delete("/", h.UserSelfDelete)

		userGroup.Use(AdminOnlyMiddleware)

		userGroup.Post("/", h.UserCreate)
		userGroup.Delete("/force/:id", h.UserForceDelete)
	}
	{
		fileGroup := r.Group("/file")

		fileGroup.Use(h.JwtMiddleware)

		fileGroup.Post("/", h.FileUpload)
		//ileGroup.Delete("/:id")

		//fileGroup.Use(AdminOnlyMiddleware)
		//fileGroup.Get("/:id")
	}
}

func (h *Handler) Data() data.Store {
	return h.data
}

func (h *Handler) Config() config.Config {
	return h.config
}

func (h *Handler) AWS() *session.Session {
	return h.aws
}
