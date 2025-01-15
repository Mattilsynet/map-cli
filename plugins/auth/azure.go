package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

var (
	clientID    string
	tenantID    string
	azureScopes []string = []string{"https://graph.microsoft.com/.default"}
)

func azureAuth() (azcore.TokenCredential, error) {
	clientID = "stuff"
	tenantID = "stuff"

	options := azidentity.DeviceCodeCredentialOptions{
		TenantID: "9e5b7d0e-770b-49e3-90ec-464fe313bdf4",
		ClientID: clientID,
		UserPrompt: func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
			fmt.Printf("Go to %s and enter the code: %s\n", message.VerificationURL, message.UserCode)
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
