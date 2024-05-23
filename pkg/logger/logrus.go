package logger

import (
	"nymphicus-service/config"
	"os"

	"github.com/sirupsen/logrus"
)

// Logger methods interface
type Logger interface {
	InitLogger()
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Panic(args ...interface{})
	Panicf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}

// Logger
type apiLogger struct {
	cfg    *config.Config
	logger *logrus.Logger
}

// App Logger constructor
func NewApiLogger(cfg *config.Config) *apiLogger {
	return &apiLogger{cfg: cfg}
}

// For mapping config logger to app logger levels
var loggerLevelMap = map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
	"panic": logrus.PanicLevel,
	"fatal": logrus.FatalLevel,
}

func (l *apiLogger) getLoggerLevel(cfg *config.Config) logrus.Level {
	level, exist := loggerLevelMap[cfg.Logger.Level]
	if !exist {
		return logrus.DebugLevel
	}

	return level
}

// Init logger
func (l *apiLogger) InitLogger() {
	logLevel := l.getLoggerLevel(l.cfg)

	l.logger = logrus.New()
	l.logger.SetOutput(os.Stderr)
	l.logger.SetLevel(logLevel)

	if l.cfg.Logger.Encoding == "console" {
		l.logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		l.logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}
}

// Logger methods

func (l *apiLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *apiLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args...)
}

func (l *apiLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *apiLogger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}

func (l *apiLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *apiLogger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}

func (l *apiLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *apiLogger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
}

func (l *apiLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *apiLogger) Panicf(template string, args ...interface{}) {
	l.logger.Panicf(template, args...)
}

func (l *apiLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *apiLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args...)
}
