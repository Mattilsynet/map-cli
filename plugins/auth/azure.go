package main

import (
	"fmt"
	"os"

	"github.com/Mattilsynet/map-cli/pkg/auth/azureauth"
	"github.com/spf13/cobra"
)

var (
	clientID string
	tenantID string
)

var azureCmd = &cobra.Command{
	Use:   "azure",
	Short: "Authenticate with device code flow",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var azureCmdLogin = &cobra.Command{
	Use:   "login",
	Short: "Login with device code flow",
	Run: func(cmd *cobra.Command, args []string) {
		auth, err := azureauth.New(
			azureauth.WithTenantID(tenantID),
			azureauth.WithClientID(clientID),
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		loginErr := auth.Login()
		if loginErr != nil {
			fmt.Println(loginErr)
			os.Exit(1)
		}
		fmt.Println("token type:", auth.TokenType())
		fmt.Println("idtoken:", auth.IDToken())
		fmt.Println("accesstoken:", auth.AccessToken())
		fmt.Println("token expires in:", auth.ExpiresIn())
	},
}

func init() {
	// Using env vars as default values here, this should probably come from viper config instead.
	azureCmd.PersistentFlags().StringVar(&clientID, "az-client-id", os.Getenv("AZURE_CLIENT_ID"), "Azure client ID")
	azureCmd.PersistentFlags().StringVar(&tenantID, "az-tenant-id", os.Getenv("AZURE_TENANT_ID"), "Azure tenant ID")

	rootCmd.AddCommand(azureCmd)
	azureCmd.AddCommand(azureCmdLogin)
	azureCmd.TraverseChildren = true
}
