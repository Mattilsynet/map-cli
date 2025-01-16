package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	_ "github.com/Mattilsynet/map-cli/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	MC_CONFIG_NAME = "config.toml"
)

var mcConfigFile string

func init() {
	pflag.StringVarP(&mcConfigFile, "config", "c", MC_CONFIG_NAME, "file to read configuration from")
}

func main() {
	// Parses all flags and makes them available in pflag.CommandLine.
	pflag.Parse()
	slog.Debug("Logger initialized")

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
