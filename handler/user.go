package handler

import (
	"github.com/RevinB/mira_server/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func (h *Handler) UserResetSecret(c *fiber.Ctx) error {
	userid := c.FormValue("userid")

	if userid == "" {
		return c.Status(fiber.StatusBadRequest).SendString("no form key 'userid' found")
	}

	userData, err := h.Data().User().GetById(userid)
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).SendString("user not found")
	} else if err != nil {
		return err
	}

	userData.SecretKey, err = utils.GenerateRandomString(32)
	if err != nil {
		return err
	}

	err = h.Data().User().Update(userData)
	if err != nil {
		return err
	}

	claims := jwt.MapClaims{
		"secret_key": userData.SecretKey,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	finalToken, err := token.SignedString(h.Config().JWTSecret)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"token": finalToken})
}
