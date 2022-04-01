package handler

import (
	"fmt"
	"github.com/RevinB/mira_server/data/file"
	"github.com/RevinB/mira_server/data/user"
	"github.com/RevinB/mira_server/utils"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gofiber/fiber/v2"
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

	fileKey, err := utils.GenerateRandomString(8)
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
	matches := rxp.FindAllString(mpFile.Filename, -1)
	fileExt = matches[len(matches)-1]

	newFileName := fileKey + fileExt

	cType := "application/octet-stream"
	if t := mpFile.Header.Get("content-type"); t != "" {
		cType = t
	}

	open, err := mpFile.Open()
	if err != nil {
		return err
	}

	//hash := md5.New()
	//if err != nil {
	//	return err
	//}
	//
	//_, err = io.Copy(hash, open)
	//if err != nil {
	//	return err
	//}

	//strHash := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	uploader := s3manager.NewUploader(h.AWS())

	upParams := &s3manager.UploadInput{
		Body:   open,
		Bucket: utils.NewStringPointer(h.Config().S3BucketName),
		//ContentMD5: utils.NewStringPointer(strHash),
		ContentType: utils.NewStringPointer(cType),
		Key:         utils.NewStringPointer(newFileName),
	}

	_, err = uploader.Upload(upParams)
	if err != nil {
		return err
	}

	dbEntry := file.Model{
		ID:            fileKey,
		FileExtension: strings.TrimPrefix(fileExt, "."),
		MIMEType:      cType,
		Owner:         userData.ID,
	}

	err = h.Data().Files().Create(&dbEntry)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).SendString(fmt.Sprintf("%s/%s", h.Config().FinalUrlBase, newFileName))
}
