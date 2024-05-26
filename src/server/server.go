package server

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"nymphicus-service/config"
	"nymphicus-service/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	ctxTimeout = 5
)

// Server struct
type Server struct {
	cfg    *config.Config
	logger logger.Logger
	mongo  *mongo.Client
	srv    *fasthttp.Server
	redis  *redis.Client
}

func (s *Server) loggingMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		next(ctx)
		end := time.Now()
		s.logger.Infof("Method: %s, URI: %s, Status: %d, Duration: %s",
			string(ctx.Method()), ctx.URI().String(), ctx.Response.StatusCode(), end.Sub(start))
	}
}

// NewServer New Server constructor
func NewServer(cfg *config.Config, logger logger.Logger, mongo *mongo.Client, redis *redis.Client) *Server {
	server := &Server{
		cfg:    cfg,
		logger: logger,
		mongo:  mongo,
		srv: &fasthttp.Server{
			Name:         "FastHTTP Server",
			ReadTimeout:  time.Second * cfg.Server.ReadTimeout,
			WriteTimeout: time.Second * cfg.Server.WriteTimeout,
		},
		redis: redis,
	}
	return server
}

func (s *Server) Run() error {
	s.srv.Handler = s.loggingMiddleware(s.handler)

	go func() {
		s.logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
		if err := s.srv.ListenAndServe(s.cfg.Server.Port); err != nil && !errors.Is(err, errors.New("http: Server closed")) {
			s.logger.Fatalf("Error starting Server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	s.logger.Info("Server Exited Properly")
	return s.srv.ShutdownWithContext(ctx)
}
