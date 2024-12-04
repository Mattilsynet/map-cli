package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "mc",
		Short: "Main command (mc) for managing tasks",
	}

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
