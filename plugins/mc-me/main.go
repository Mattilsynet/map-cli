// mc-me/main.go
package main

import (
	"fmt"

	"github.com/Mattilsynet/map-cli/internal/config"
	"github.com/Mattilsynet/map-cli/plugins/mc-me/handler"
	"github.com/spf13/cobra"
)

func main() {
	nc, err := config.CurrentConfig.Nats.GetConnection()
	if err != nil {
		fmt.Printf("Error connecting to NATS: %v\n", err)
		return
	}
	handler := handler.New(nc)
	rootCmd := &cobra.Command{
		Use:     "me",
		Short:   "Managed Environment (me) plugin",
		Aliases: []string{"me"},
	}
	rootCmd.AddCommand(&cobra.Command{
		Use:   "apply",
		Short: "Create or update a managed-environment",
		Run: func(cmd *cobra.Command, args []string) {
			err := handler.HandleCobraCommand(cmd, args)
			fmt.Printf("Error: %v\n", err)
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "create",
		Short: "Create a managed-environment",
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "update",
		Short: "Update a managed-environment",
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "Delete a managed-environment",
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get a managed-environment",
		Run: func(cmd *cobra.Command, args []string) {
			handler.HandleCobraCommand(cmd, args)
		},
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
