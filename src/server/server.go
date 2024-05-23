package server

import (
	"context"
	"errors"
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
	srv    *fasthttp.Server
}

// NewServer New Server constructor
func NewServer(cfg *config.Config, logger logger.Logger) *Server {
	server := &Server{
		cfg:    cfg,
		logger: logger,
		srv: &fasthttp.Server{
			Name:         "FastHTTP Server",
			ReadTimeout:  time.Second * cfg.Server.ReadTimeout,
			WriteTimeout: time.Second * cfg.Server.WriteTimeout,
		},
	}
	return server
}

func (s *Server) Run() error {
	s.srv.Handler = s.handler

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
