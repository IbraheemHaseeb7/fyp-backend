package controllers

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
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
		return c.File(fmt.Sprintf("%s/%s/%s/%s",
			filepath.Join("storage/private"), c.Param("level1"), c.Param("level2"), c.Param("filename")))
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

		id := c.Get("auth_user_id")

		hash := utils.GenerateHash(map[string]string{
			"id": fmt.Sprintf("%.f", id.(float64)),
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
		os.Mkdir(fmt.Sprintf("%s/%s/%.f", filepath.Join("storage"), medium, id.(float64)), 0777)
		os.Mkdir(fmt.Sprintf("%s/%s/%.f/%s", filepath.Join("storage"), medium, id.(float64), folder), 0777)

		// storing file in the folder system
		outFilePath := filepath.Join("storage/"+medium+"/"+fmt.Sprintf("%.f", id.(float64))+"/"+folder+"/", hash+"."+fileExt)
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
			"filepath": fmt.Sprintf("%s/%.f/%s/%s.%s", medium, id, folder, hash, fileExt),
		}

		cr.APIResponse.Status = "Successfully uploaded file"
		cr.APIResponse.StatusCode = 201
		return cr.SendSuccessResponse(&c)
	}
}

func DeleteFiles(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenHeader := strings.Split(c.Request().Header.Get("Authorization"), " ")
		if len(tokenHeader) < 1 {
			cr.APIResponse.Status = "Could not authorize"
			cr.APIResponse.StatusCode = 401
			cr.APIResponse.Error = "Token not found in the header"
			return cr.SendErrorResponse(&c)
		}

		uuid := watermill.NewUUID()
		utils.Requests.Store(uuid, make(chan pubsub.PubsubMessage))

		cr.Publisher.PublishMessage(pubsub.PubsubMessage{
			Payload: map[string]string{
				"token": tokenHeader[1],
			},
			Entity:    "files",
			Operation: "GET_CLAIMS",
			UUID:      uuid,
			Topic:     "img->auth",
		})
		authResp := <-utils.Requests.Load(uuid)
		utils.Requests.Delete(uuid)

		payload, ok := authResp.Payload.(map[string]any)
		if !ok {
			return cr.SendErrorResponse(&c)
		}

		if status := payload["verified"]; status == false {
			cr.APIResponse.Status = "Token invalid"
			cr.APIResponse.StatusCode = 401
			return cr.SendErrorResponse(&c)
		}

		userId := payload["id"]
		type Request struct {
			Files []string `json:"files"`
		}
		var reqBody Request
		if err := cr.BindAndValidate(&reqBody, &c); err != nil {
			return cr.SendErrorResponse(&c)
		}

		for i := range reqBody.Files {
			filePath := strings.Split(reqBody.Files[i], "/")
			filePath[1] = fmt.Sprintf("%.f", userId)
			newFilePath:= strings.Join(filePath, "/")
			wd, _ := os.Getwd()
			err := os.Remove(wd + "/storage/" + newFilePath) 
			if err != nil {
				fmt.Println(wd)
				continue
			}
		}

		cr.APIResponse.Data = "Successfully deleted all files"
		return cr.SendSuccessResponse(&c)
	}
}

func KeepOnlyFiles(cr ControllerRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Get("auth_user_id")

		userId := fmt.Sprintf("%.f", id.(float64))

		type Request struct {
			Files []string `json:"files"`
			Dir string `json:"dir"`
		}
		var reqBody Request
		if err := cr.BindAndValidate(&reqBody, &c); err != nil {
			return cr.SendErrorResponse(&c)
		}

		// getting working directory
		wd, err := os.Getwd()
		if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		// creating temp dir if does not exist
		destDir := filepath.Join(wd, "storage", reqBody.Dir, "temp")
		err = os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		// Moving all the files in a temp directory
		for _, path := range reqBody.Files {
			split := strings.Split(path, "/")
			if len(split) == 4 && split[1] == userId {
				sourcePath := filepath.Join(wd, "storage", path)
				destPath := filepath.Join(destDir, split[3])
				err = os.Rename(sourcePath, destPath)
				if err != nil {
					cr.APIResponse.Error = err.Error()
					return cr.SendErrorResponse(&c)
				}
			}
		}

		// Removing all the files in the dir
		entries, err := os.ReadDir(filepath.Join(wd, "storage", reqBody.Dir))
		if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			err := os.Remove(filepath.Join(wd, "storage", reqBody.Dir, entry.Name()))
			if err != nil {
				cr.APIResponse.Error = err.Error()
				return cr.SendErrorResponse(&c)
			}
		}

		// Moving all the files out of the temp dir
		tempEntries, err := os.ReadDir(filepath.Join(wd, "storage", reqBody.Dir, "temp"))
		if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}
		for _, entry := range tempEntries {
			if entry.IsDir() {
				continue
			}

			sourcePath := filepath.Join(wd, "storage", reqBody.Dir, "temp", entry.Name())
			destPath := filepath.Join(wd, "storage", reqBody.Dir, entry.Name())

			err := os.Rename(sourcePath, destPath)
			if err != nil {
				cr.APIResponse.Error = err.Error()
				return cr.SendErrorResponse(&c)
			}
		}

		// removing the temp dir
		err = os.Remove(path.Join(wd, "storage", reqBody.Dir, "temp"))
		if err != nil {
			cr.APIResponse.Error = err.Error()
			return cr.SendErrorResponse(&c)
		}

		cr.APIResponse.Data = "Kept specified files and deleted the rest"
		return cr.SendSuccessResponse(&c)
	}
}

