package server

import (
	"github.com/valyala/fasthttp"
	"nymphicus-service/src/controllers"
	"nymphicus-service/src/controllers/health"
	controllerv2 "nymphicus-service/src/controllers/v2"
	"nymphicus-service/src/repository"
	service "nymphicus-service/src/services"
)

func (s *Server) handler(ctx *fasthttp.RequestCtx) {

	indexRepository := controllers.NewController(s.cfg, s.logger, s.mongo)
	sessionRepository := repository.NewSessionRepository(s.mongo)

	videoService := service.NewVideoService(s.cfg)

	checkRecordingController := controllers.NewCheckRecordingController(s.cfg, s.logger, s.redis)
	writeVideoDataController := controllerv2.NewWriteVideoDataController(s.cfg, s.logger, sessionRepository, videoService)

	switch string(ctx.Path()) {
	case "/write":
		indexRepository.ControllerSDK(ctx)
	case "/v2/write":
		writeVideoDataController.WriteVideoData(ctx)
	case "/check-recording":
		checkRecordingController.ValidateAccessKey(ctx)
	case "/health":
		health.CheckHandler(ctx)
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}
