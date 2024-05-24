package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"io"
	"log"
	"mime/multipart"
	"nymphicus-service/config"
	"nymphicus-service/pkg/logger"
	"nymphicus-service/pkg/utils"
	"nymphicus-service/src/models"
	"path/filepath"
)

type Controller interface {
	ControllerSDK(ctx *fasthttp.RequestCtx)
}

type controller struct {
	config *config.Config
	logger logger.Logger
}

func NewController(
	config *config.Config,
	logger logger.Logger,
) Controller {
	return &controller{
		config: config,
		logger: logger,
	}
}

func (c *controller) ControllerSDK(ctx *fasthttp.RequestCtx) {
	multipartForm, err := ctx.MultipartForm()
	if err != nil {
		utils.HandleRequestError(ctx, err, c.logger)
		return
	}

	fileHeader, err := extractFile(multipartForm)
	if err != nil {
		utils.HandleRequestError(ctx, err, c.logger)
		return
	}

	actions, err := extractActionsData(multipartForm)
	if err != nil {
		utils.HandleRequestError(ctx, err, c.logger)
		return
	}

	timeLines := utils.GetTimeLines(actions.Activities)

	go func() {
		if err := sendToPythonEndpoint(fileHeader, timeLines, uuid.New().String()); err != nil {
			log.Printf("Failed to send data to Python endpoint: %v", err)
		}
	}()

	device, err := extractDeviceData(multipartForm)
	if err != nil {
		utils.HandleRequestError(ctx, err, c.logger)
		return
	}

	response := fmt.Sprintf("File received: %s\nDevice: %+v\nActions: %+v", fileHeader.Filename, device, actions)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(response))
}

func extractFile(form *multipart.Form) (*multipart.FileHeader, error) {
	files := form.File["file"]
	if len(files) == 0 {
		return nil, fmt.Errorf("file is missing")
	}
	return files[0], nil
}

func extractDeviceData(form *multipart.Form) (models.Device, error) {
	var device models.Device
	deviceData := form.Value["device"]
	if len(deviceData) == 0 {
		return device, fmt.Errorf("device data is missing")
	}
	if err := json.Unmarshal([]byte(deviceData[0]), &device); err != nil {
		return device, fmt.Errorf("failed to parse device data")
	}
	return device, nil
}

func extractActionsData(form *multipart.Form) (models.Action, error) {
	var actions models.Action
	actionsData := form.Value["actions"]
	if len(actionsData) == 0 {
		return actions, fmt.Errorf("actions data is missing")
	}
	if err := json.Unmarshal([]byte(actionsData[0]), &actions); err != nil {
		return actions, fmt.Errorf("failed to parse actions data")
	}
	return actions, nil
}

func sendToPythonEndpoint(fileHeader *multipart.FileHeader, timeLines []utils.TimeLine, sessionId string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	part, err := writer.CreateFormFile("file", filepath.Base(fileHeader.Filename))
	if err != nil {
		return err
	}
	if _, err = io.Copy(part, file); err != nil {
		return err
	}

	if err = addFormField(writer, "timeLines", timeLines); err != nil {
		return err
	}

	if err = writer.WriteField("sessionId", sessionId); err != nil {
		return err
	}

	if err = writer.Close(); err != nil {
		return err
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod("POST")
	req.Header.SetContentType(writer.FormDataContentType())
	req.SetRequestURI("http://10.0.0.106:8000/send_binary_data")
	req.SetBody(body.Bytes())

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	if err := fasthttp.Do(req, resp); err != nil {
		return err
	}

	log.Printf("Response status code: %d\n", resp.StatusCode())
	log.Printf("Response body: %s\n", resp.Body())

	return nil
}

func addFormField(writer *multipart.Writer, fieldName string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return writer.WriteField(fieldName, string(jsonData))
}
