package handler

import (
	"github.com/RevinB/mira_server/data/file"
	"github.com/RevinB/mira_server/data/user"
	"github.com/RevinB/mira_server/utils"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func (h *Handler) UserResetSecret(c *fiber.Ctx) error {
	userid := c.Params("id")

	if userid == "" {
		return c.Status(fiber.StatusBadRequest).SendString("no form key 'userid' found")
	}

	userData, err := h.Data.Users.GetById(userid)
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).SendString("user not found")
	} else if err != nil {
		return err
	}

	userData.SecretKey, err = utils.GenerateRandomString(32)
	if err != nil {
		return err
	}

	err = h.Data.Users.Update(userData)
	if err != nil {
		return err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"secret_key": userData.SecretKey,
	})

	finalToken, err := token.SignedString(h.Config.JWTSecret)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).SendString(finalToken)
}

func (h *Handler) UserCreate(c *fiber.Ctx) error {
	dbEntry := user.Model{
		IsAdmin: false,
	}

	err := h.Data.Users.Create(&dbEntry)
	if err != nil {
		return err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"secret_key": dbEntry.SecretKey,
	})

	finalToken, err := token.SignedString(h.Config.JWTSecret)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).SendString(finalToken)
}

func deleteAllUserFiles(h *Handler, files []file.Model, userData *user.Model) error {
	objectIds := make([]*s3.ObjectIdentifier, len(files))
	cfPaths := make([]*string, len(files))
	for i, j := range files {
		fullName := j.ID + "." + j.FileExtension
		cfPaths[i] = utils.NewStringPointer("/" + fullName)
		objectIds[i] = &s3.ObjectIdentifier{
			Key: utils.NewStringPointer(fullName),
		}
	}

	doi := s3.DeleteObjectsInput{
		Bucket: h.Config.S3BucketName,
		Delete: &s3.Delete{
			Objects: objectIds,
		},
	}

	inva := cloudfront.CreateInvalidationInput{
		DistributionId: h.Config.CloudfrontDistID,
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: utils.NewStringPointer(strconv.Itoa(int(time.Now().UnixNano()))),
			Paths: &cloudfront.Paths{
				Items:    cfPaths,
				Quantity: utils.NewInt64Pointer(int64(len(cfPaths))),
			},
		},
	}

	_, err := h.S3Client.DeleteObjects(&doi)
	if err != nil {
		return err
	}

	_, err = h.CfClient.CreateInvalidation(&inva)
	if err != nil {
		return err
	}

	err = h.Data.Files.DeleteAllFromUser(userData.ID)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) UserSelfDelete(c *fiber.Ctx) error {
	userData := c.Locals("user").(*user.Model)

	files, err := h.Data.Files.GetAllByUser(userData.ID)
	if err != nil {
		return err
	}

	// avoid unnecessary aws calls
	if len(files) > 0 {
		err = deleteAllUserFiles(h, files, userData)
		if err != nil {
			return err
		}
	}

	err = h.Data.Users.Delete(userData)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

// in admin-only middleware
func (h *Handler) UserForceDelete(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return fiber.ErrBadRequest
	}

	userData, err := h.Data.Users.GetById(id)
	if err == gorm.ErrRecordNotFound {
		return fiber.ErrNotFound
	} else if err != nil {
		return err
	}

	files, err := h.Data.Files.GetAllByUser(userData.ID)
	if err != nil {
		return err
	}

	if len(files) > 0 {
		err = deleteAllUserFiles(h, files, userData)
		if err != nil {
			return err
		}
	}

	err = h.Data.Users.Delete(userData)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
