package controllers

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/valyala/fasthttp"
	"nymphicus-service/config"
	"nymphicus-service/pkg/logger"
	"nymphicus-service/pkg/utils"
)

type CheckRecordingController interface {
	ValidateAccessKey(ctx *fasthttp.RequestCtx)
}

type checkRecordingController struct {
	config      *config.Config
	logger      logger.Logger
	redisClient *redis.Client
}

func NewCheckRecordingController(
	config *config.Config,
	logger logger.Logger,
	redisClient *redis.Client,
) CheckRecordingController {
	return &checkRecordingController{
		config:      config,
		logger:      logger,
		redisClient: redisClient,
	}
}

func (c *checkRecordingController) ValidateAccessKey(ctx *fasthttp.RequestCtx) {
	accessKey := string(ctx.QueryArgs().Peek("key"))
	if len(accessKey) == 0 {
		utils.HandleRequestError(ctx, errors.New("missing 'key' query parameter"), c.logger)
		return
	}

	exists, err := c.redisClient.Exists(context.Background(), accessKey).Result()
	if err != nil {
		utils.HandleRequestError(ctx, err, c.logger)
		return
	}

	if exists == 1 {
		ctx.SetStatusCode(fasthttp.StatusOK)
	} else {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
	}
}
