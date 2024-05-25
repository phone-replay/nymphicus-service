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

	mongoClient, err := database.ConnectionDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, cfg.Server.SSL)

	s := server.NewServer(cfg, appLogger, mongoClient)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
