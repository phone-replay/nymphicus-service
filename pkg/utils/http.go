package utils

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"nymphicus-service/pkg/httpErrors"
	"nymphicus-service/pkg/logger"
)

// GetRequestID retrieves the X-Request-ID header from the request.
func GetRequestID(ctx *fasthttp.RequestCtx) string {
	return string(ctx.Request.Header.Peek("X-Request-ID"))
}

// GetIPAddress retrieves the client IP address from the request context.
func GetIPAddress(ctx *fasthttp.RequestCtx) string {
	return ctx.RemoteIP().String()
}

// LogResponseError logs an error with the request ID and client IP address.
func LogResponseError(ctx *fasthttp.RequestCtx, logger logger.Logger, err error) {
	logger.Errorf(
		"ErrResponseWithLog, RequestID: %s, IPAddress: %s, Error: %s",
		GetRequestID(ctx),
		GetIPAddress(ctx),
		err,
	)
}

// HandleRequestError handles an error by logging it and sending an error response.
func HandleRequestError(ctx *fasthttp.RequestCtx, err error, logger logger.Logger) {
	status, body := httpErrors.ErrorResponse(err)
	LogResponseError(ctx, logger, err)
	ctx.SetStatusCode(status)
	ctx.SetContentType("application/json")
	_ = json.NewEncoder(ctx.Response.BodyWriter()).Encode(body)
}

func GetConfigPath(configPath string) string {
	if configPath == "production" {
		return "./config/config-production"
	}
	return "./config/config-local"
}
