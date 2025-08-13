package zitadel

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Mattilsynet/map-cli/internal/config"
	"github.com/spf13/viper"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
)

const (
	MapCliClientId = "332875394606759939" // mattilsynet
	// MapCliClientId = "333004134489849859" // home office
	// MapCliClientId = "332861563872542723"
	Issuer = "http://localhost:8080"
)

// INFO: code snippets from: https://github.com/zitadel/oidc/blob/main/example/client/device/device.go
// INFO: Guide followed: https://zitadel.com/docs/guides/integrate/login/oidc/device-authorization
func DeviceLogin() error {
	// TODO: issuer and clientid from env not inline
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
	defer stop()
	clientID := MapCliClientId
	issuer := Issuer
	scopess := "email profile openid"
	scopes := strings.Split(scopess, " ")
	clientSecret := ""
	// might need to use cookies in the future, not sure.
	provider, err := rp.NewRelyingPartyOIDC(ctx, issuer, clientID, clientSecret, "", scopes)
	if err != nil {
		fmt.Println("error creating provider %s", err.Error())
		panic("error creating zitadel provider")
	}

	fmt.Println("starting device authorization flow")
	resp, err := rp.DeviceAuthorization(ctx, scopes, provider, nil)
	if err != nil {
		fmt.Println("error authorizing towards zitadel with: ", err)
		panic(err)
	}
	fmt.Println("resp", resp)
	fmt.Printf("\nPlease browse to %s and enter code %s\n", resp.VerificationURI, resp.UserCode)

	fmt.Println("start polling")
	token, err := rp.DeviceAccessToken(ctx, resp.DeviceCode, time.Duration(resp.Interval)*time.Second, provider)
	if err != nil {
		panic(err)
	}
	fmt.Printf("successfully obtained token: %#v", token)
	return execLogin(token.AccessToken, token.IDToken, token.ExpiresIn)
}

func execLogin(zBearerToken, zIdToken string, expiresIn uint64) error {
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("error reading in viper config", err)
		return err
	}
	configInFile := config.CurrentConfig
	err = viper.Unmarshal(configInFile)
	if err != nil {
		fmt.Printf("error unmarshalling viper config", err)
		return err
	}
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second).Unix()
	viper.Set("zitadel.BearerToken", zBearerToken)
	viper.Set("zitadel.IdToken", zIdToken)
	viper.Set("zitadel.ExpiresAt", expiresAt)
	err = viper.WriteConfig()
	if err != nil {
		fmt.Printf("error writing viper config", err)
		return err
	}
	return nil
}
