package main

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
	"nymphicus-service/src/server"
	"os"
	"path/filepath"
	"sort"
)

type DeviceInfo struct {
	BatteryLevel     float64 `json:"batteryLevel"`
	Brand            string  `json:"brand"`
	CurrentNetwork   string  `json:"currentNetwork"`
	Device           string  `json:"device"`
	InstallID        string  `json:"installID"`
	Language         string  `json:"language"`
	Manufacturer     string  `json:"manufacturer"`
	Model            string  `json:"model"`
	OsVersion        string  `json:"osVersion"`
	Platform         string  `json:"platform"`
	ScreenResolution string  `json:"screenResolution"`
	SdkVersion       int     `json:"sdkVersion"`
	SessionId        string  `json:"sessionId"`
	TotalRAM         string  `json:"totalRAM"`
	TotalStorage     string  `json:"totalStorage"`
}

type Action struct {
	Activities []Activity `json:"activities"`
}

type Activity struct {
	ActivityName string    `json:"activityName"`
	Gestures     []Gesture `json:"gestures"`
	ID           string    `json:"id"`
}

type Gesture struct {
	ActivityId  string `json:"activityId"`
	Coordinates string `json:"coordinates"`
	CreatedAt   string `json:"createdAt"`
	GestureType string `json:"gestureType"`
	TargetTime  string `json:"targetTime"`
}

func getTimeLines(activities []Activity) []TimeLine {
	var timeLines []TimeLine
	for _, activity := range activities {
		for _, gesture := range activity.Gestures {
			timeLine := TimeLine{
				Coordinates: gesture.Coordinates,
				GestureType: gesture.GestureType,
				TargetTime:  gesture.TargetTime,
			}
			timeLines = append(timeLines, timeLine)
		}
	}
	sort.Slice(timeLines, func(i, j int) bool {
		return timeLines[i].TargetTime < timeLines[j].TargetTime
	})

	return timeLines
}

func uploadHandler(ctx *fasthttp.RequestCtx) {
	// Extrair boundary do header Content-Type
	multipartForm, err := ctx.MultipartForm()
	if err != nil {
		ctx.Error("failed to parse multipart form", fasthttp.StatusBadRequest)
		return
	}

	// Extrair arquivo
	files := multipartForm.File["file"]
	if len(files) == 0 {
		ctx.Error("file is missing", fasthttp.StatusBadRequest)
		return
	}
	fileHeader := files[0] // Isto é um *multipart.FileHeader

	// Extrair JSON do dispositivo
	deviceData := multipartForm.Value["device"]
	if len(deviceData) == 0 {
		ctx.Error("device data is missing", fasthttp.StatusBadRequest)
		return
	}
	var device DeviceInfo
	if err := json.Unmarshal([]byte(deviceData[0]), &device); err != nil {
		ctx.Error("failed to parse device data", fasthttp.StatusBadRequest)
		return
	}

	// Extrair JSON das ações
	actionsData := multipartForm.Value["actions"]
	if len(actionsData) == 0 {
		ctx.Error("actions data is missing", fasthttp.StatusBadRequest)
		return
	}
	var actions Action
	if err := json.Unmarshal([]byte(actionsData[0]), &actions); err != nil {
		ctx.Error("failed to parse actions data", fasthttp.StatusBadRequest)
		return
	}

	timeLines := getTimeLines(actions.Activities)

	go func() {
		err := sendToPythonEndpoint(fileHeader, timeLines, uuid.New().String())
		if err != nil {
			log.Printf("failed to send data to Python endpoint: %v", err)
		}
	}()

	// Resposta de sucesso
	response := fmt.Sprintf("File received: %s\nDevice: %+v\nActions: %+v", fileHeader.Filename, device, actions)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(response))
}

type TimeLine struct {
	Coordinates string `json:"coordinates"`
	GestureType string `json:"gestureType"`
	TargetTime  string `json:"targetTime"`
}

func sendToPythonEndpoint(fileHeader *multipart.FileHeader, timeLines []TimeLine, sessionId string) error {
	// Buffer para o corpo da solicitação
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Adicionar o arquivo
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(fileHeader.Filename))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	// Serializar timeLines como JSON
	timeLinesJSON, err := json.Marshal(timeLines)
	if err != nil {
		return err
	}

	// Adicionar JSON de timeLines
	err = writer.WriteField("timeLines", string(timeLinesJSON))
	if err != nil {
		return err
	}

	// Adicionar sessionId
	err = writer.WriteField("sessionId", sessionId)
	if err != nil {
		return err
	}

	// Fechar o writer para finalizar a montagem do corpo da solicitação
	err = writer.Close()
	if err != nil {
		return err
	}

	// Criar a solicitação usando fasthttp
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod("POST")
	req.Header.SetContentType(writer.FormDataContentType())
	req.SetRequestURI("http://10.0.0.106:8000/send_binary_data")
	req.SetBody(body.Bytes())

	// Executar a solicitação
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err = fasthttp.Do(req, resp)
	if err != nil {
		return err
	}

	// Opcionalmente, processar a resposta aqui
	fmt.Printf("Response status code: %d\n", resp.StatusCode())
	fmt.Printf("Response body: %s\n", resp.Body())

	return nil
}

func main() {
	configPath := utils.GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, cfg.Server.SSL)

	s := server.NewServer(cfg, appLogger)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
