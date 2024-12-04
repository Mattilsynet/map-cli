// mc-me/main.go
package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

/*
TODO: implement call towards map.<map-type>.session.id, somehow generate a session from logged in nats session
*/
func main() {
	rootCmd := &cobra.Command{
		Use:     "mc-me",
		Short:   "Managed Environment (me) plugin",
		Aliases: []string{"me"},
	}
	rootCmd.AddCommand(&cobra.Command{
		Use:   "apply",
		Short: "Create or update a managed-environment",
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

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
