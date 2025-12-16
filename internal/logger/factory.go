package logger

import (
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	AppName     string
	Environment string
	LogLevel    string
	LogFile     string
	EnableJSON  bool
	EnableColor bool
}

func NewLogger(cfg Config) zerolog.Logger {
	level := parseLevel(cfg.LogLevel)
	zerolog.SetGlobalLevel(level)

	var writers []io.Writer

	if cfg.Environment != "production" || cfg.EnableColor {
		writers = append(writers, createConsoleWriter(cfg))
	}

	if cfg.LogFile != "" {
		writers = append(writers, createFileWriter(cfg.LogFile))
	}

	if cfg.Environment == "production" {
		writers = append(writers, createSyslogWriter())
	}

	multiWriter := io.MultiWriter(writers...)

	logger := zerolog.New(multiWriter).
		With().
		Timestamp().
		Caller().
		Str("app", cfg.AppName).
		Str("env", cfg.Environment).
		Logger()

	log.Logger = logger

	return logger
}

func parseLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

func createConsoleWriter(cfg Config) io.Writer {
	if cfg.EnableJSON {
		return os.Stderr
	}

	return zerolog.ConsoleWriter{
		Out:        os.Stderr,
		NoColor:    !cfg.EnableColor,
		TimeFormat: "2006-01-02 15:04:05",
	}
}

func createFileWriter(filename string) io.Writer {
	file, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	return file
}

func createSyslogWriter() io.Writer {
	return nil
}
