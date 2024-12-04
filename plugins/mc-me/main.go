// mc-me/main.go
package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "me",
		Short:   "Managed Environment (me) plugin",
		Aliases: []string{"me"},
	}
	rootCmd.AddCommand(&cobra.Command{
		Use:   "apply",
		Short: "Create or update a managed-environment",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("apply called")
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
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
