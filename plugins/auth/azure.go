package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	clientID    string
	tenantID    string
	azureScopes []string = []string{"https://graph.microsoft.com/.default"}
)

var azureCmd = &cobra.Command{
	Use:   "azure",
	Short: "Authenticate with device code flow",
}

var azureCmdLogin = &cobra.Command{
	Use:   "login",
	Short: "Login with device code flow",
	Run: func(cmd *cobra.Command, args []string) {
		cred, err := azureAuth()
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		fmt.Println(cred)
	},
}

func init() {
	// Using env vars as default values here, this should probably come from viper config instead.
	pflag.StringVar(&clientID, "az-client-id", os.Getenv("AZURE_CLIENT_ID"), "Azure client ID")
	pflag.StringVar(&tenantID, "az-tenant-id", os.Getenv("AZURE_TENANT_ID"), "Azure tenant ID")
}

func azureAuth() (azcore.TokenCredential, error) {
	options := azidentity.DeviceCodeCredentialOptions{
		TenantID: tenantID,
		ClientID: clientID,
		UserPrompt: func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
			fmt.Printf("%s", message.Message)
			return nil
		},
	}

	// Create a DeviceCodeCredential
	cred, err := azidentity.NewDeviceCodeCredential(&options)
	if err != nil {
		return nil, fmt.Errorf("failed to create device code credential: %w", err)
	}

	// Acquire a token for the specified scopes
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	token, err := cred.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: azureScopes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to acquire token: %w", err)
	}

	// Print the acquired token for verification (for demo purposes only)
	fmt.Printf("Access Token: %s\n", token.Token)

	return cred, nil
}
