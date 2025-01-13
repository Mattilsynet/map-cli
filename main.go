package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	MC_CONFIG_NAME = "config.toml"
)

var (
	mcLogLevel   slog.Level
	mcConfigFile string
)

func init() {
	pflag.StringVar(&mcConfigFile, "config", "config.toml", "file to read configuration from")
	var mcLogLevelInput string
	pflag.StringVar(&mcLogLevelInput, "log-level", "info", "log level (debug, info, warn, error)")

	initLogger(mcLogLevelInput)
}

func initLogger(level string) {
	logLevel, err := parseLogLevel(level)
	if err != nil {
		fmt.Println("Invalid log level provided: " + err.Error())
		os.Exit(1)
	}
	logger := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	})
	slog.SetDefault(slog.New(logger))
}

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

func main() {
	rootCmd := &cobra.Command{
		Use:   "mc",
		Short: "Main command (mc) for managing tasks",
	}

	rootCmd.Flags().AddFlagSet(pflag.CommandLine)

	rootCmd.AddCommand(&cobra.Command{
		Use:     "managed-environment",
		Short:   "Managed Environment (me)",
		Aliases: []string{"me"},
		Run: func(cmd *cobra.Command, args []string) {
			err := execPlugin("mc-me", args...)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to execute plugin: %v\n", err)
				os.Exit(1)
			}
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:     "auth",
		Short:   "Authentication plugin",
		Aliases: []string{"a"},
		Run: func(cmd *cobra.Command, args []string) {
			err := execPlugin("mc-auth", args...)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to execute plugin: %v\n", err)
				os.Exit(1)
			}
		},
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func execPlugin(pluginName string, args ...string) error {
	path, err := exec.LookPath(pluginName)
	if err != nil {
		return fmt.Errorf("plugin '%s' not found in PATH", pluginName)
	}

	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
