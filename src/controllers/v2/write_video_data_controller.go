package controller_v2

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"log"
	"mime/multipart"
	"nymphicus-service/config"
	"nymphicus-service/enum"
	"nymphicus-service/pkg/logger"
	"nymphicus-service/pkg/utils"
	"nymphicus-service/src"
	"nymphicus-service/src/models"
	"nymphicus-service/src/repository"
	service "nymphicus-service/src/services"
	"strconv"
	"time"
)

type WriteVideoDataController interface {
	WriteVideoData(ctx *fasthttp.RequestCtx)
}

type writeVideoData struct {
	config            *config.Config
	logger            logger.Logger
	sessionRepository repository.SessionRepository
	videoService      service.VideoService
}

func NewWriteVideoDataController(
	config *config.Config,
	logger logger.Logger,
	sessionRepository repository.SessionRepository,
	videoService service.VideoService,
) WriteVideoDataController {
	return &writeVideoData{
		config:            config,
		logger:            logger,
		sessionRepository: sessionRepository,
		videoService:      videoService,
	}
}

func (c *writeVideoData) WriteVideoData(ctx *fasthttp.RequestCtx) {
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

	activityGestureLogs, err := extractActivityGestureLogs(multipartForm)
	if err != nil {
		utils.HandleRequestError(ctx, err, c.logger)
		return
	}

	duration, err := extractDurationData(multipartForm)
	if err != nil {
		utils.HandleRequestError(ctx, err, c.logger)
		return
	}

	session, err := createSession(key, device, duration)
	if err != nil {
		utils.HandleRequestError(ctx, err, c.logger)
		return
	}

	err = c.sessionRepository.SaveActionsToMongo(session)
	if err != nil {
		utils.HandleRequestError(ctx, err, c.logger)
		return
	}

	go func() {
		if err := c.videoService.RequestGenerateVideo(fileHeader, activityGestureLogs, session.ID, strconv.FormatInt(duration, 10)); err != nil {
			if err := c.sessionRepository.UpdateSessionStatusToError(key); err != nil {
				log.Printf("Failed to update session status to error: %v", err)
			}
			log.Printf("Failed to generate video: %v", err)
		}
	}()

	response := fmt.Sprintf("File received: %s\nDevice: %+v\nActivity Gesture Logs: %+v\nDuration: %d", fileHeader.Filename, device, activityGestureLogs, duration)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(response))
}

// createSession creates a new session object.
func createSession(key string, device models.Device, duration int64) (models.Session, error) {
	return models.Session{
		ID:        uuid.New().String(),
		Key:       key,
		Device:    device,
		Status:    enum.InProgress.String(),
		CreatedAt: time.Now(),
		Duration:  duration,
	}, nil
}

// extractActivityGestureLogs extracts ActivityGestureLogs from the multipart form.
func extractActivityGestureLogs(form *multipart.Form) (src.ActivityGestureLogs, error) {
	var activityGestureLogs src.ActivityGestureLogs
	fileHeader, err := getFileHeader(form, "activityGestureLogs")
	if err != nil {
		return activityGestureLogs, err
	}

	data, err := utils.ReadGzipFile(fileHeader)
	if err != nil {
		return activityGestureLogs, err
	}

	if err := json.Unmarshal(data, &activityGestureLogs.Activities); err != nil {
		return activityGestureLogs, fmt.Errorf("failed to parse activityGestureLogs data: %v", err)
	}

	return activityGestureLogs, nil
}

// extractFile extracts the file from the multipart form.
func extractFile(form *multipart.Form) (*multipart.FileHeader, error) {
	return getFileHeader(form, "file")
}

// extractDeviceData extracts device data from the multipart form.
func extractDeviceData(form *multipart.Form) (models.Device, error) {
	var device models.Device
	deviceData, err := getFormValue(form, "device")
	if err != nil {
		return device, err
	}
	if err := json.Unmarshal([]byte(deviceData), &device); err != nil {
		return device, fmt.Errorf("failed to parse device data")
	}
	return device, nil
}

// extractDurationData extracts duration data from the multipart form.
func extractDurationData(form *multipart.Form) (int64, error) {
	durationData, err := getFormValue(form, "duration")
	if err != nil {
		return 0, err
	}
	duration, err := strconv.ParseInt(durationData, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration data")
	}
	return duration, nil
}

// getFileHeader retrieves a file header from the multipart form.
func getFileHeader(form *multipart.Form, key string) (*multipart.FileHeader, error) {
	files := form.File[key]
	if len(files) == 0 {
		return nil, fmt.Errorf("%s file is missing", key)
	}
	return files[0], nil
}

// getFormValue retrieves a form value from the multipart form.
func getFormValue(form *multipart.Form, key string) (string, error) {
	values := form.Value[key]
	if len(values) == 0 {
		return "", fmt.Errorf("%s data is missing", key)
	}
	return values[0], nil
}
