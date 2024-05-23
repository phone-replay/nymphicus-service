package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
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

// GetConfigPath returns the configuration path based on the provided configPath string.
func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker"
	}
	return "./config/config-local"
}

// ReadRequest reads and validates a JSON request.
func ReadRequest(ctx *fasthttp.RequestCtx, request interface{}) error {
	if err := json.Unmarshal(ctx.PostBody(), request); err != nil {
		return httpErrors.NewBadRequestError(err.Error())
	}

	validate := validator.New()

	if err := validate.Struct(request); err != nil {
		var vErrs validator.ValidationErrors
		if errors.As(err, &vErrs) {
			var validationErrors []string
			for _, vErr := range vErrs {
				validationErrors = append(validationErrors, fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", vErr.Field(), vErr.Tag()))
			}
			return httpErrors.NewBadRequestError(validationErrors)
		}
		return httpErrors.NewBadRequestError(err.Error())
	}

	return nil
}
