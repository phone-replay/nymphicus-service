package main

import (
	"log"
	"nymphicus-service/config"
	"nymphicus-service/database"
	"nymphicus-service/pkg/logger"
	"nymphicus-service/pkg/utils"
	"nymphicus-service/src/server"
	"os"
)

func main() {
	configPath := utils.GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, cfg.Server.SSL)

	redisClient, err := database.NewRedisClient(cfg)
	if err != nil {
		appLogger.Fatalf("Redis init: %s", err)
	}
	appLogger.Infof("Redis connected")

	mongoClient, err := database.ConnectionDatabase(cfg, appLogger)
	if err != nil {
		log.Fatal(err)
	}

	s := server.NewServer(cfg, appLogger, mongoClient, redisClient)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
