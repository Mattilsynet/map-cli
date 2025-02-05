/*
	This should be moved elsewhere but, need to doc this somewhere.

	On a Mac, execute the following to find the intune mdm client certificate:

	security find-certificate -a -c "IntuneMDM" -p > client-cert.pem
	security export -k -t priv -p -c "IntuneMDM" -o private_key.pem


*/

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log/slog"
	"net/http"
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

var azureCmdLoginCustom = &cobra.Command{
	Use:   "login2",
	Short: "Login with device code flow manually",
	Run: func(cmd *cobra.Command, args []string) {
		// Step 1: Request device code
		deviceCodeResp, err := requestDeviceCode()
		if err != nil {
			fmt.Println("Failed to request device code:", err)
			return
		}

		fmt.Println(deviceCodeResp.Message) // Show message to user

		// Step 2: Poll for token response
		tokenResp, err := pollForToken(deviceCodeResp.DeviceCode)
		if err != nil {
			fmt.Println("Failed to obtain token:", err)
			return
		}

		// Step 3: Extract Access & ID Token
		fmt.Println("Access Token:", tokenResp.AccessToken)
		fmt.Println("ID Token:", tokenResp.IDToken)
	},
}

func init() {
	// Using env vars as default values here, this should probably come from viper config instead.
	azureCmd.PersistentFlags().StringVar(&clientID, "az-client-id", os.Getenv("AZURE_CLIENT_ID"), "Azure client ID")
	azureCmd.PersistentFlags().StringVar(&tenantID, "az-tenant-id", os.Getenv("AZURE_TENANT_ID"), "Azure tenant ID")

	rootCmd.AddCommand(azureCmd)
	azureCmd.AddCommand(azureCmdLogin)
	azureCmd.AddCommand(azureCmdLoginCustom)
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

type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
	Message         string `json:"message"`
}

func requestDeviceCode() (*DeviceCodeResponse, error) {
	data := "client_id=" + clientID + "&scope=" + "openid profile email"

	authURL := "https://login.microsoftonline.com/" + tenantID + "/oauth2/v2.0"

	req, err := http.NewRequest("POST", authURL+"/devicecode", bytes.NewBufferString(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result DeviceCodeResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"` // Extract ID Token
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

func pollForToken(deviceCode string) (*TokenResponse, error) {
	data := fmt.Sprintf(
		"grant_type=urn:ietf:params:oauth:grant-type:device_code&device_code=%s&client_id=%s",
		deviceCode, clientID,
	)
	authURL := "https://login.microsoftonline.com/" + tenantID + "/oauth2/v2.0"
	client := &http.Client{}
	for {
		req, _ := http.NewRequest("POST", authURL+"/token", bytes.NewBufferString(data))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == 200 {
			var result TokenResponse
			err = json.Unmarshal(body, &result)
			if err != nil {
				return nil, err
			}
			return &result, nil
		}

		// If authorization is still pending, wait before retrying
		time.Sleep(5 * time.Second)
	}
}
