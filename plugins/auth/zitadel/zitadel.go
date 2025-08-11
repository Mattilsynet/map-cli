package zitadel

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/zitadel/oidc/v3/pkg/client/rp"
)

const (
	MapCliClientId = "332875394606759939"
	// MapCliClientId = "332861563872542723"
	Issuer = "http://localhost:8080"
)

// INFO: code snippets from: https://github.com/zitadel/oidc/blob/main/example/client/device/device.go
// INFO: Guide followed: https://zitadel.com/docs/guides/integrate/login/oidc/device-authorization
func DeviceLogin() {
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
}
