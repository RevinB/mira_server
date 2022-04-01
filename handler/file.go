package handler

import (
	"crypto/sha256"
	"fmt"
	"github.com/RevinB/mira_server/data/file"
	"github.com/RevinB/mira_server/data/user"
	"github.com/RevinB/mira_server/utils"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gofiber/fiber/v2"
	"io"
	"regexp"
	"strconv"
	"strings"
)

const maxUploadLen = 1000 * 1000 * 5000

func (h *Handler) FileUpload(c *fiber.Ctx) error {
	userData := c.Locals("user").(*user.Model)

	contLen := c.Get("Content-Length")
	if contLen == "" {
		return fiber.ErrLengthRequired
	}

	iContLen, err := strconv.Atoi(contLen)
	if err != nil {
		return err
	}

	if iContLen > maxUploadLen {
		return fiber.ErrRequestEntityTooLarge
	}

	mpFile, err := c.FormFile("file")
	if err != nil {
		return fiber.ErrBadRequest
	}

	newFileName, err := utils.GenerateRandomString(8)
	if err != nil {
		return err
	}

	// get file extension
	fileExt := ""
	var rxp *regexp.Regexp
	if strings.HasSuffix(mpFile.Filename, ".gz") {
		rxp = regexp.MustCompile(`(\.[A-Za-z0-9-_]+\.gz)$`)
	} else {
		rxp = regexp.MustCompile(`(\.[A-Za-z0-9-_]+)$`)
	}

	// get last match
	matches := rxp.FindAllString(fileExt, -1)
	fileExt = matches[len(matches)-1]

	cType := "application/octet-stream"
	if t := c.Get("Content-Type"); t != "" {
		cType = t
	}

	open, err := mpFile.Open()
	if err != nil {
		return err
	}

	hash := sha256.New()
	_, err = io.Copy(hash, open)
	if err != nil {
		return err
	}

	strHash := string(hash.Sum(nil))

	uploader := s3manager.NewUploader(h.AWS())

	upParams := &s3manager.UploadInput{
		Body:           open,
		Bucket:         utils.NewStringPointer(h.Config().S3BucketName),
		ChecksumSHA256: utils.NewStringPointer(strHash),
		Key:            utils.NewStringPointer(newFileName + fileExt),
		Metadata: map[string]*string{
			"Content-Type": utils.NewStringPointer(cType),
		},
	}

	_, err = uploader.Upload(upParams)
	if err != nil {
		return err
	}

	dbEntry := file.Model{
		ID:            newFileName,
		FileExtension: fileExt,
		MIMEType:      cType,
		Owner:         userData.ID,
	}

	err = h.Data().Files().Create(&dbEntry)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).SendString(fmt.Sprintf("%s/%s.%s", h.Config().FinalUrlBase, newFileName, fileExt))
}
