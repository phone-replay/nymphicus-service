package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"log"
	"mime/multipart"
	"nymphicus-service/config"
	"nymphicus-service/enum"
	"nymphicus-service/pkg/logger"
	"nymphicus-service/pkg/utils"
	"nymphicus-service/src/models"
	"path/filepath"
	"time"
)

type Controller interface {
	ControllerSDK(ctx *fasthttp.RequestCtx)
}

type controller struct {
	config      *config.Config
	logger      logger.Logger
	mongoClient *mongo.Client
}

func NewController(
	config *config.Config,
	logger logger.Logger,
	mongoClient *mongo.Client,
) Controller {
	return &controller{
		config:      config,
		logger:      logger,
		mongoClient: mongoClient,
	}
}

func (c *controller) ControllerSDK(ctx *fasthttp.RequestCtx) {
	key := string(ctx.QueryArgs().Peek("key"))
	if len(key) == 0 {
		utils.HandleRequestError(ctx, errors.New("missing 'key' query parameter"), c.logger)
		return
	}

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

	device, err := extractDeviceData(multipartForm)
	if err != nil {
		utils.HandleRequestError(ctx, err, c.logger)
		return
	}

	session, err := extractActionsData(multipartForm)
	if err != nil {
		utils.HandleRequestError(ctx, err, c.logger)
		return
	}

	session.Device = device
	session.Key = key
	sessionId := uuid.New().String()
	session.ID = sessionId
	session.Status = enum.InProgress.String()
	session.CreatedAt = time.Now()

	err = c.saveActionsToMongo(session)
	if err != nil {
		utils.HandleRequestError(ctx, err, c.logger)
		return
	}

	timeLines := utils.GetTimeLines(session.Activities)

	go func() {
		if err := generateVideo(fileHeader, timeLines, sessionId, c.config); err != nil {
			err := c.updateSessionStatusToError(key)
			if err != nil {
				utils.HandleRequestError(ctx, err, c.logger)
				return
			}
			log.Printf("Failed to send data to Python endpoint: %v", err)
		}
	}()

	response := fmt.Sprintf("File received: %s\nDevice: %+v\nActions: %+v", fileHeader.Filename, device, session)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(response))
}

func (c *controller) saveActionsToMongo(actions models.Session) error {
	collection := c.mongoClient.Database("mongo_db").Collection("sessions")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, actions)
	return err
}

func (c *controller) updateSessionStatusToError(key string) error {
	collection := c.mongoClient.Database("mongo_db").Collection("sessions")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"key": key}
	update := bson.M{
		"$set": bson.M{
			"status": enum.Error.String(),
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
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

func extractActionsData(form *multipart.Form) (models.Session, error) {
	var session models.Session
	actionsData := form.Value["actions"]
	if len(actionsData) == 0 {
		return session, fmt.Errorf("session data is missing")
	}
	if err := json.Unmarshal([]byte(actionsData[0]), &session); err != nil {
		return session, fmt.Errorf("failed to parse session data")
	}
	return session, nil
}

func generateVideo(fileHeader *multipart.FileHeader, timeLines []utils.TimeLine, sessionId string, config *config.Config) error {
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
	req.SetRequestURI(config.Services.OtididaeURL)
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
