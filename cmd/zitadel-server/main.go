// INFO: This is a standalone test server. It's ment to be used in conjunction with mc auth zitadel login to validate and test out the token and test the authorization zitadel server on localhost:8080
// INFO: Read more on device code flow: https://zitadel.com/docs/guides/integrate/login/oidc/device-authorization
// INFO: And on how to setup your local zitadel server:https://zitadel.com/docs/self-hosting/deploy/compose
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/oidcclient"
	"github.com/zitadel/oidc/v3/pkg/oidcclient/httphelper"
)

func main() {
	ctx := context.Background()

	issuer := "http://localhost:8080" // Local Zitadel
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		log.Fatalf("failed to get provider: %v", err)
	}

	keySet := oidc.NewRemoteKeySet(ctx, provider.Endpoint().JWKSURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}
		rawToken := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate token
		idToken, err := oidc.ParseToken(rawToken, keySet)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid token: %v", err), http.StatusUnauthorized)
			return
		}
		fmt.Println(idToken)
		// Optional: Validate claims (audience, issuer)

		// If we get here, token is valid
		fmt.Fprintln(w, "hello world")
	})

	fmt.Println("Server running at http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}
