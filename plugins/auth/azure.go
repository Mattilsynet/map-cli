/*
	This should be moved elsewhere but, need to doc this somewhere.

	On a Mac, execute the following to find the intune mdm client certificate:

	security find-certificate -a -c "IntuneMDM" -p > client-cert.pem
	security export -k -t priv -p -c "IntuneMDM" -o private_key.pem


*/

package main

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/spf13/cobra"
)

var (
	clientID       string
	tenantID       string
	clientCertPath string
)

var (
	azureScopes           []string = []string{"https://graph.microsoft.com/.default"}
	azureManagementScopes []string = []string{"https://management.azure.com/.default"}
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
		cred, err := azureDeviceCodeFlowAuth()
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		fmt.Println(cred)
	},
}

var azureCmdCertLogin = &cobra.Command{
	Use:   "cert",
	Short: "Login with Intune client credentials",
	Run: func(cmd *cobra.Command, args []string) {
		if len(clientCertPath) == 0 {
			_ = cmd.Help()
			return
		}
		token, err := azureClientCertificateCredential()
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		fmt.Println(token)
	},
}

func init() {
	// Using env vars as default values here, this should probably come from viper config instead.
	azureCmd.PersistentFlags().StringVar(&clientID, "az-client-id", os.Getenv("AZURE_CLIENT_ID"), "Azure client ID")
	azureCmd.PersistentFlags().StringVar(&tenantID, "az-tenant-id", os.Getenv("AZURE_TENANT_ID"), "Azure tenant ID")
	azureCmdCertLogin.PersistentFlags().StringVar(&clientCertPath, "az-client-cert-path", "", "Path to Intune MDM Client Certificate (PEM format)")
	azureCmdCertLogin.MarkFlagRequired("az-client-cert-path")

	rootCmd.AddCommand(azureCmd)
	azureCmd.AddCommand(azureCmdLogin)
	azureCmd.AddCommand(azureCmdCertLogin)
	azureCmd.TraverseChildren = true
}

func azureDeviceCodeFlowAuth() (azcore.TokenCredential, error) {
	options := azidentity.DeviceCodeCredentialOptions{
		TenantID: tenantID,
		ClientID: clientID,
		UserPrompt: func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
			fmt.Printf("%s", message.Message)
			return nil
		},
		ClientOptions: azcore.ClientOptions{},
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

func azureClientCertificateCredential() (*azcore.AccessToken, error) {
	certData, err := os.ReadFile(clientCertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate: %w", err)
	}

	// Decode the PEM-encoded certificate
	block, _ := pem.Decode(certData)
	if block == nil {
		return nil, errors.New("failed to parse PEM block from certificate")
	}

	// Parse the X.509 certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Create a ClientCertificateCredential
	cred, err := azidentity.NewClientCertificateCredential(tenantID, clientID, []*x509.Certificate{cert}, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ClientCertificateCredential: %w", err)
	}

	// Retrieve an access token
	token, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: azureManagementScopes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to acquire token: %w", err)
	}

	return &token, nil
}
