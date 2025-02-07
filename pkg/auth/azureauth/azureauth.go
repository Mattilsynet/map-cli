package azureauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
	Message         string `json:"message"`
}

type Auth struct {
	TenantID   string
	ClientID   string
	token      *TokenResponse
	deviceCode *DeviceCodeResponse
}

type AuthOption func(*Auth)

func New(options ...AuthOption) (*Auth, error) {
	azureAuth := &Auth{}
	for _, option := range options {
		option(azureAuth)
	}
	if !azureAuth.valid() {
		return nil, fmt.Errorf("invalid azure auth configuration")
	}
	return azureAuth, nil
}

func (a *Auth) valid() bool {
	validations := []bool{
		a.TenantID != "",
		a.ClientID != "",
	}

	for _, v := range validations {
		if !v {
			return false
		}
	}
	return true
}

func WithTenantID(tenantID string) AuthOption {
	return func(a *Auth) {
		a.TenantID = tenantID
	}
}

func WithClientID(clientID string) AuthOption {
	return func(a *Auth) {
		a.ClientID = clientID
	}
}

func (a *Auth) AuthURL() string {
	return "https://login.microsoftonline.com/" + a.TenantID + "/oauth2/v2.0"
}

func (a *Auth) Login() error {
	err := a.requestDeviceCode()
	if err != nil {
		return fmt.Errorf("Failed to request device code:", err)
	}

	fmt.Println(a.deviceCode.Message)

	// Step 2: Poll for token response
	pollErr := a.pollForToken()
	if pollErr != nil {
		return fmt.Errorf("Failed to obtain token:", pollErr)
	}
	return nil
}

func (a *Auth) IDToken() string {
	return a.token.IDToken
}

func (a *Auth) AccessToken() string {
	return a.token.AccessToken
}

func (a *Auth) requestDeviceCode() error {
	data := "client_id=" + a.ClientID + "&scope=" + "openid profile email"

	req, err := http.NewRequest("POST", a.AuthURL()+"/devicecode", bytes.NewBufferString(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &a.deviceCode)
	if err != nil {
		return err
	}

	return nil
}

func (a *Auth) pollForToken() error {
	data := fmt.Sprintf(
		"grant_type=urn:ietf:params:oauth:grant-type:device_code&device_code=%s&client_id=%s",
		a.deviceCode.DeviceCode, a.ClientID,
	)
	client := &http.Client{}
	for {
		req, _ := http.NewRequest("POST", a.AuthURL()+"/token", bytes.NewBufferString(data))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 200 {
			err = json.Unmarshal(body, a.token)
			if err != nil {
				return err
			}
			return nil
		}

		// If authorization is still pending, wait before retrying
		time.Sleep(5 * time.Second)
	}
}
