package server

import (
	"github.com/valyala/fasthttp"
	"nymphicus-service/src/controllers"
)

// handler is the FastHTTP request handler
func (s *Server) handler(ctx *fasthttp.RequestCtx) {
	// Define your request handling logic here

	c := controllers.NewController(s.cfg, s.logger)

	switch string(ctx.Path()) {
	case "/upload":
		c.ControllerSDK(ctx)
		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		ctx.Response.SetBodyString("Hello, World!")
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}
