package server

import "github.com/valyala/fasthttp"

// handler is the FastHTTP request handler
func (s *Server) handler(ctx *fasthttp.RequestCtx) {
	// Define your request handling logic here
	switch string(ctx.Path()) {
	case "/":
		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		ctx.Response.SetBodyString("Hello, World!")
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}
