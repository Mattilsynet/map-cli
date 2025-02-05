package main

import (
	"fmt"
	"os"

	"github.com/Mattilsynet/map-cli/plugins/component/pkg/config"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "component",
		Short:   "Component plugin",
		Aliases: []string{"c"},
	}

	generate := &cobra.Command{
		Use:     "generate",
		Short:   "Generate a WasmCloud component",
		Aliases: []string{"gen", "g"},
		Run: func(cmd *cobra.Command, args []string) {
			componentConfig := promptUser()
			// cue validate schema
			componentProject := config.CreateComponentProject(componentConfig)
			config.GenerateFilesPrompt(componentProject)
		},
	}
	rootCmd.AddCommand(generate)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func promptUser() string {
	// component-name
	// root-path
	// nats-core
	// nats-jetstream
	// nats-kv
	// license
	panic("unimplemented")
}
