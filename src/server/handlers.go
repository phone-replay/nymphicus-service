package server

import (
	"github.com/valyala/fasthttp"
	"nymphicus-service/src/controllers"
)

// handler is the FastHTTP request handler
func (s *Server) handler(ctx *fasthttp.RequestCtx) {
	indexRepository := controllers.NewController(s.cfg, s.logger, s.mongo)
	checkRecordingController := controllers.NewCheckRecordingController(s.cfg, s.logger, s.redis)
	switch string(ctx.Path()) {
	case "/write":
		indexRepository.ControllerSDK(ctx)
	case "/check-recording":
		checkRecordingController.ValidateAccessKey(ctx)
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}
