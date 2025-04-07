package controllers

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/IbraheemHaseeb7/fyp-backend/http/services"
	"github.com/IbraheemHaseeb7/fyp-backend/utils"
	"github.com/IbraheemHaseeb7/pubsub"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/labstack/echo/v4"
)

type FileHandler struct {
	F *services.FilesServices
	Cr *ControllerRequest
}

func ReadOnePrivateFile(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {

		tokenHeader := strings.Split(c.Request().Header.Get("Authorization"), " ")
		if len(tokenHeader) > 1 {
			uuid := watermill.NewUUID()
			utils.Requests[uuid] = make(chan pubsub.PubsubMessage)

			cr.Publisher.PublishMessage(pubsub.PubsubMessage{
				Payload: map[string]string{
					"token": tokenHeader[1],
				},
				Entity:    "files",
				Operation: "VERIFY_TOKEN",
				UUID:      uuid,
				Topic:     "img->auth",
			})
			authResp := <-utils.Requests[uuid]
			delete(utils.Requests, uuid)

			payload, ok := authResp.Payload.(map[string]any)
			if !ok {
				return cr.SendErrorResponse(&c)
			}

			if status := payload["verified"]; status == false {
				cr.APIResponse.Status = "Token invalid"
				cr.APIResponse.StatusCode = 401
				return cr.SendErrorResponse(&c)
			}
			return c.File(fmt.Sprintf("%s/%s/%s/%s",
				filepath.Join("storage/private"), c.Param("level1"), c.Param("level2"), c.Param("filename")))
		}

		cr.APIResponse.Status = "Token not found"
		cr.APIResponse.StatusCode = 401
		return cr.SendErrorResponse(&c)
	}
}

func (f *FileHandler) ReadOnePublicFile(c echo.Context) error {
	return c.File(fmt.Sprintf("%s/%s/%s/%s",
		filepath.Join("storage/public"), c.Param("level1"), c.Param("level2"), c.Param("filename")))
}

func UploadOneFile(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {

		// limiting the file size to 10MB
		c.Request().ParseMultipartForm(10 << 20)

		folder := c.Request().FormValue("folder")
		if folder == "" {
			cr.APIResponse.Status = "Validation Failed"
			cr.APIResponse.StatusCode = 400
			cr.APIResponse.Error = "Did not provide folder name"
			return cr.SendErrorResponse(&c)
		}

		medium := c.Request().FormValue("medium")
		if medium == "" {
			cr.APIResponse.Status = "Validation Failed"
			cr.APIResponse.StatusCode = 400
			cr.APIResponse.Error = "Did not provide medium"
			return cr.SendErrorResponse(&c)
		}

		file, fileHeader, err := c.Request().FormFile("file")
		if err != nil {
			cr.APIResponse.Status = "Could not parse file"
			cr.APIResponse.StatusCode = 500
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		defer file.Close()

		tokenHeader := strings.Split(c.Request().Header.Get("Authorization"), " ")
		if len(tokenHeader) < 1 {
			cr.APIResponse.Status = "Could not authorize"
			cr.APIResponse.StatusCode = 401
			cr.APIResponse.Error = "Token not found in the header"
			return cr.SendErrorResponse(&c)
		}

		uuid := watermill.NewUUID()
		utils.Requests[uuid] = make(chan pubsub.PubsubMessage)

		cr.Publisher.PublishMessage(pubsub.PubsubMessage{
			Payload: map[string]string{
				"token": tokenHeader[1],
			},
			Entity:    "files",
			Operation: "GET_CLAIMS",
			UUID:      uuid,
			Topic:     "img->auth",
		})
		authResp := <-utils.Requests[uuid]
		delete(utils.Requests, uuid)

		payload, ok := authResp.Payload.(map[string]any)
		if !ok {
			return cr.SendErrorResponse(&c)
		}

		if status := payload["verified"]; status == false {
			cr.APIResponse.Status = "Token invalid"
			cr.APIResponse.StatusCode = 401
			return cr.SendErrorResponse(&c)
		}

		hash := utils.GenerateHash(map[string]string{
			"id": fmt.Sprintf("%.f", payload["id"].(float64)),
			"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
		})

		// extracting file extension
		fileExt := strings.Split(fileHeader.Filename, ".")[1]

		var img image.Image
		switch fileExt {
		case "png", "PNG":
			img, err = png.Decode(file)
			if err != nil {
				cr.APIResponse.Status = "Could not parse file"
				cr.APIResponse.StatusCode = 500
				cr.APIResponse.Error = err.Error()
				return cr.SendErrorResponse(&c)
			}
		case "jpg", "jpeg", "JPG", "JPEG":
			img, err = jpeg.Decode(file)
			if err != nil {
				cr.APIResponse.Status = "Could not parse file"
				cr.APIResponse.StatusCode = 500
				cr.APIResponse.Error = err.Error()
				return cr.SendErrorResponse(&c)
			}
		default:
			cr.APIResponse.Status = "Could not parse file"
			cr.APIResponse.StatusCode = 500
			cr.APIResponse.Error = "File format not supported here."
			return cr.SendErrorResponse(&c)
		}

		// creating directory if doesn't exist
		os.Mkdir(fmt.Sprintf("%s/%s/%.f", filepath.Join("storage"), medium, payload["id"].(float64)), 0777)
		os.Mkdir(fmt.Sprintf("%s/%s/%.f/%s", filepath.Join("storage"), medium, payload["id"].(float64), folder), 0777)

		// storing file in the folder system
		outFilePath := filepath.Join("storage/"+medium+"/"+fmt.Sprintf("%.f", payload["id"].(float64))+"/"+folder+"/", hash+"."+fileExt)
		outFile, err := os.Create(outFilePath)
		if err != nil {
			cr.APIResponse.Status = "Could not parse file"
			cr.APIResponse.StatusCode = 500
			cr.APIResponse.Error = err.Error()
			return cr.SendResponse(&c)
		}
		defer outFile.Close()

		// Encode the image back to PNG format
		switch fileExt {
		case "png", "PNG":
			if err := png.Encode(outFile, img); err != nil {
				cr.APIResponse.Status = "Could not parse file"
				cr.APIResponse.StatusCode = 500
				cr.APIResponse.Error = err.Error()
				return cr.SendErrorResponse(&c)
			}
		case "jpg", "jpeg", "JPG", "JPEG":
			if err := jpeg.Encode(outFile, img, &jpeg.Options{Quality: 75}); err != nil {
				cr.APIResponse.Status = "Could not parse file"
				cr.APIResponse.StatusCode = 500
				cr.APIResponse.Error = err.Error()
				return cr.SendErrorResponse(&c)
			}
		default:
			cr.APIResponse.Status = "Could not parse file"
			cr.APIResponse.StatusCode = 500
			cr.APIResponse.Error = "File format not supported here."
			return cr.SendErrorResponse(&c)
		}

		cr.APIResponse.Data = map[string]string{
			"filepath": fmt.Sprintf("%s/%.f/%s/%s.%s", medium, payload["id"], folder, hash, fileExt),
		}
		return cr.SendSuccessResponse(&c)
	}
}

func DeleteFiles(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		return cr.SendSuccessResponse(&c)
	}
}
