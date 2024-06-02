package server

import (
	"github.com/valyala/fasthttp"
	"nymphicus-service/src/controllers"
)

func HealthCheckHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody([]byte(`{"status":"ok"}`))
}

func (s *Server) handler(ctx *fasthttp.RequestCtx) {
	indexRepository := controllers.NewController(s.cfg, s.logger, s.mongo)
	checkRecordingController := controllers.NewCheckRecordingController(s.cfg, s.logger, s.redis)
	switch string(ctx.Path()) {
	case "/write":
		indexRepository.ControllerSDK(ctx)
	case "/check-recording":
		checkRecordingController.ValidateAccessKey(ctx)
	case "/health":
		HealthCheckHandler(ctx)
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}
