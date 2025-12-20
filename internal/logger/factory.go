package logger

import (
	"io"
	"os"
	"strings"

	"github.com/ExplosiveGM/wasted/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewLogger(cfg *config.Config) zerolog.Logger {
	level := parseLevel(cfg.Log.Level)
	zerolog.SetGlobalLevel(level)

	var writers []io.Writer

	if cfg.App.Env != "production" || cfg.Log.EnableColor {
		writers = append(writers, createConsoleWriter(cfg))
	}

	if cfg.Log.File != "" {
		writers = append(writers, createFileWriter(cfg.Log.File))
	}

	if cfg.App.Env == "production" {
		writers = append(writers, createSyslogWriter())
	}

	multiWriter := io.MultiWriter(writers...)

	logger := zerolog.New(multiWriter).
		With().
		Timestamp().
		Caller().
		Str("app", cfg.App.Name).
		Str("env", cfg.App.Env).
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

func createConsoleWriter(cfg *config.Config) io.Writer {
	if cfg.Log.EnableJson {
		return os.Stderr
	}

	return zerolog.ConsoleWriter{
		Out:        os.Stderr,
		NoColor:    !cfg.Log.EnableColor,
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
