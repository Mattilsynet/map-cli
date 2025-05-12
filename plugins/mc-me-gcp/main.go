// mc-me/main.go
package main

import (
	"fmt"
	"time"

	"github.com/Mattilsynet/map-cli/internal/config"
	"github.com/Mattilsynet/map-cli/plugins/mc-me-gcp/handler"
	"github.com/spf13/cobra"
)

func main() {
	// TODO: We need to make sure that we're not reliant on nats connection to actually show available commands.
	nc, err := config.CurrentConfig.Nats.GetConnection()
	if err != nil {
		fmt.Printf("Error connecting to NATS: %v\n", err)
		return
	}
	handler := handler.New(nc)
	rootCmd := &cobra.Command{
		Use:     "me",
		Short:   "Managed Environment for gcp (me-gcp) plugin",
		Aliases: []string{"me"},
	}
	rootCmd.AddCommand(&cobra.Command{
		Use:   "apply",
		Short: "Create or update a managed-environment in gcp",
		Run: func(cmd *cobra.Command, args []string) {
			err := handler.HandleCobraCommand(cmd, args)
			if err != nil {
				// TODO: We should provide normal human-interaction errors and not just print the error
				fmt.Println("Error: ", err)
			}
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "Delete a managed-environment in gcp",
		Run: func(cmd *cobra.Command, args []string) {
			err := handler.HandleCobraCommand(cmd, args)
			if err != nil {
				// TODO: We should provide normal human-interaction errors and not just print the error
				fmt.Println("Error: ", err)
			}
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get a managed-environment in gcp",
		Run: func(cmd *cobra.Command, args []string) {
			err := handler.HandleCobraCommand(cmd, args)
			if err != nil {
				fmt.Println("Error: ", err)
			}
		},
	})
	time.Sleep(2000)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
