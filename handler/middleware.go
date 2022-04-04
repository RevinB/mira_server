package handler

import (
	"github.com/RevinB/mira_server/data/user"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"strings"
)

func (h *Handler) JwtMiddleware(c *fiber.Ctx) error {
	authHeader := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")

	token, err := jwt.Parse(authHeader, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrExpectationFailed
		}

		return h.Config.JWTSecret, nil
	})
	if err != nil {
		return fiber.ErrUnauthorized
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		secKey := claims["secret_key"].(string)
		userData, err := h.Data.Users.GetByKey(secKey)
		if err != nil {
			return err
		}

		c.Locals("user", userData)

		return c.Next()
	} else {
		return fiber.ErrForbidden
	}
}

func AdminOnlyMiddleware(c *fiber.Ctx) error {
	if !c.Locals("user").(*user.Model).IsAdmin {
		return fiber.ErrForbidden
	}
	return c.Next()
}
