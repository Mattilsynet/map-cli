// mc-me/main.go
package main

import (
	"fmt"

	"github.com/Mattilsynet/map-cli/internal/config"
	"github.com/Mattilsynet/map-cli/plugins/mc-org/handler"
	"github.com/spf13/cobra"
)

func main() {
	// TODO: refactor validation outside of this module, in config init could be a great start
	if config.CurrentConfig == nil {
		fmt.Println("Could not load config. Please run 'mc auth zitadel login'")
		return

	}
	configs := config.CurrentConfig
	_ = configs
	bearerToken := config.CurrentConfig.Zitadel.BearerToken
	if bearerToken == "" {
		fmt.Println("You need to be authenticated to use this plugin. Please run 'mc auth zitadel login', no bearer token found")
		return
	}
	if config.CurrentConfig.Zitadel.IsExpired() {
		fmt.Println("Your authentication has expired. Please run 'mc auth zitadel login'")
		return
	}
	handler := handler.New(bearerToken)
	rootCmd := &cobra.Command{
		Use:   "organization",
		Short: "org",
	}
	rootCmd.AddCommand(&cobra.Command{
		Use:   "apply",
		Short: "Create or update an organization in zitadel",
		Run: func(cmd *cobra.Command, args []string) {
			err := handler.HandleCobraCommand(cmd, args)
			if err != nil {
				fmt.Println("Error: ", err)
			}
		},
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
