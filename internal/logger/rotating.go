package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

func InitWithRotation(env string) zerolog.Logger {
	var writers []io.Writer

	if env == "development" {
		writers = append(writers, getConsoleWriter())
	}

	writers = append(writers, getFileWriter(env))

	multiWriter := io.MultiWriter(writers...)

	return zerolog.New(multiWriter).
		With().
		Timestamp().
		Caller().
		Str("env", env).
		Logger()
}

func getConsoleWriter() io.Writer {
	return zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		FormatLevel: func(i interface{}) string {
			level := i.(string)
			switch level {
			case "debug":
				return "\x1b[36mDBG\x1b[0m"
			case "info":
				return "\x1b[32mINF\x1b[0m"
			case "warn":
				return "\x1b[33mWRN\x1b[0m"
			case "error":
				return "\x1b[31mERR\x1b[0m"
			default:
				return level
			}
		},
	}
}

func getFileWriter(env string) io.Writer {
	return &lumberjack.Logger{
		Filename:   getLogFilePath(env),
		MaxSize:    100,  // MB
		MaxBackups: 30,   // files
		MaxAge:     90,   // days
		Compress:   true, // using gzip for old logs
		LocalTime:  true,
	}
}

func getLogFilePath(env string) string {
	logDir := "logs"
	if env == "production" {
		logDir = "/var/log/wasted"
	}

	os.MkdirAll(logDir, 0755)

	return filepath.Join(logDir, "app.log")
}
