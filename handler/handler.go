package handler

import (
	"github.com/RevinB/mira_server/config"
	"github.com/RevinB/mira_server/data"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Data   *data.Store
	Config config.Config
	AWS    *session.Session
}

func NewHandler(data *data.Store, cfg config.Config, a *session.Session) Handler {
	return Handler{
		Data:   data,
		Config: cfg,
		AWS:    a,
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
		fileGroup.Delete("/:fileid", h.FileDelete)

		//fileGroup.Use(AdminOnlyMiddleware)
		//fileGroup.Get("/:id")
	}
}
