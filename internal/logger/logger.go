package logger

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

// Command-line flag for log level
var mcLogLevelInput string

func init() {
	// Define the log-level flag
	pflag.StringVar(&mcLogLevelInput, "log-level", "info", "log level (debug, info, warn, error)")

	// Set a default logger
	setupLogger("info")
}

// setupLogger initializes the global default logger
func setupLogger(level string) {
	logLevel, err := parseLogLevel(level)
	if err != nil {
		fmt.Println("Invalid log level provided: " + err.Error())
		os.Exit(1)
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	})
	slog.SetDefault(slog.New(handler)) // Set as the global default logger
}

// Reinitialize reinitializes the logger after flags are parsed
func Reinitialize() {
	setupLogger(mcLogLevelInput)
}

// parseLogLevel maps a string level to `slog.Level`
func parseLogLevel(level string) (slog.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, errors.New("unsupported log level")
	}
}
