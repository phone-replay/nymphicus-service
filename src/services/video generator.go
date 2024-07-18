package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"nymphicus-service/config"
	"nymphicus-service/src"
	"path/filepath"

	"github.com/valyala/fasthttp"
)

type VideoService interface {
	RequestGenerateVideo(fileHeader *multipart.FileHeader, timeLines src.ActivityGestureLogs, sessionId string, duration string) error
}

type videoService struct {
	config *config.Config
}

func NewVideoService(config *config.Config) VideoService {
	return &videoService{
		config: config,
	}
}

func (v *videoService) RequestGenerateVideo(fileHeader *multipart.FileHeader, timeLines src.ActivityGestureLogs, sessionId string, duration string) error {
	if fileHeader == nil {
		return errors.New("fileHeader cannot be nil")
	}
	if sessionId == "" {
		return errors.New("sessionId cannot be empty")
	}
	if duration == "" {
		return errors.New("duration cannot be empty")
	}

	body, contentType, err := createRequestBody(fileHeader, timeLines, sessionId, duration)
	if err != nil {
		return err
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod("POST")
	req.Header.SetContentType(contentType)

	if v.config.Services.OtididaeURL == "" {
		return errors.New("OtididaeURL cannot be empty")
	}

	req.SetRequestURI(v.config.Services.OtididaeURL)
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

func createRequestBody(fileHeader *multipart.FileHeader, timeLines src.ActivityGestureLogs, sessionId string, duration string) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := fileHeader.Open()
	if err != nil {
		return nil, "", err
	}
	defer func(file multipart.File) {
		if err := file.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(file)

	part, err := writer.CreateFormFile("file", filepath.Base(fileHeader.Filename))
	if err != nil {
		return nil, "", err
	}
	if _, err = io.Copy(part, file); err != nil {
		return nil, "", err
	}

	if err = addFormField(writer, "timeLines", timeLines); err != nil {
		return nil, "", err
	}

	if err = writer.WriteField("sessionId", sessionId); err != nil {
		return nil, "", err
	}

	if err = writer.WriteField("duration", duration); err != nil {
		return nil, "", err
	}

	if err = writer.Close(); err != nil {
		return nil, "", err
	}

	return body, writer.FormDataContentType(), nil
}

func addFormField(writer *multipart.Writer, fieldName string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return writer.WriteField(fieldName, string(jsonData))
}
