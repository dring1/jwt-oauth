package logger

import (
	"io"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type Service interface {
	Info(...interface{})
	Infof(string, ...interface{})
	Debug(...interface{})
	Debugf(string, ...interface{})
}

type Config struct {
	Environment string
	Level       string
	Endpoint    io.Writer
}

func NewLoggerService(c *Config) Service {
	var formatter log.Formatter
	formatter = new(log.TextFormatter)
	if c.Environment == "PRODUCTION" {
		formatter = new(log.JSONFormatter)
	}
	return &log.Logger{
		Out:       c.Endpoint,
		Formatter: formatter,
		Level:     getLevel(c.Level),
	}
}

func getLevel(level string) log.Level {
	l := strings.ToLower(level)
	switch l {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	}
	return log.InfoLevel
}
