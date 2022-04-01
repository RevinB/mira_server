package handler

import (
	"github.com/RevinB/mira_server/data/user"
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"secret_key": userData.SecretKey,
	})

	finalToken, err := token.SignedString(h.Config().JWTSecret)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).SendString(finalToken)
}

func (h *Handler) UserCreate(c *fiber.Ctx) error {
	dbEntry := user.Model{
		IsAdmin: false,
	}

	err := h.Data().User().Create(&dbEntry)
	if err != nil {
		return err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"secret_key": dbEntry.SecretKey,
	})

	finalToken, err := token.SignedString(h.Config().JWTSecret)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).SendString(finalToken)
}

func (h *Handler) UserSelfDelete(c *fiber.Ctx) error {
	userData := c.Locals("user").(*user.Model)

	// TODO delete files

	return h.Data().User().Delete(userData)
}

func (h *Handler) UserForceDelete(c *fiber.Ctx) error {
	id := c.FormValue("id")

	if id == "" {
		return fiber.ErrBadRequest
	}

	userData, err := h.Data().User().GetById(id)
	if err == gorm.ErrRecordNotFound {
		return fiber.ErrNotFound
	} else if err != nil {
		return err
	}

	err = h.Data().User().Delete(userData)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
