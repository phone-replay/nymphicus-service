package server

import (
	"github.com/valyala/fasthttp"
	"nymphicus-service/src/controllers"
)

// handler is the FastHTTP request handler
func (s *Server) handler(ctx *fasthttp.RequestCtx) {
	c := controllers.NewController(s.cfg, s.logger, s.mongo)

	switch string(ctx.Path()) {
	case "/write":
		c.ControllerSDK(ctx)
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}
