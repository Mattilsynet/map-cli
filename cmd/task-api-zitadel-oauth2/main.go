// INFO: This is a standalone test server. It's ment to be used in conjunction with mc auth zitadel login to validate and test out the token and test the authorization zitadel server on localhost:8080
// INFO: And on how to setup your local zitadel server:https://zitadel.com/docs/self-hosting/deploy/compose
// INFO: Guide used for the code below: https://zitadel.com/docs/examples/secure-api/go
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/exp/slog"

	"github.com/nats-io/nats.go"
	client2 "github.com/zitadel/oidc/v3/pkg/client"
	"github.com/zitadel/oidc/v3/pkg/client/rs"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/http/middleware"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

var (
	// flags to be provided for running the example server
	domain = flag.String("domain", "", "your ZITADEL instance domain (only works with localhost, since we set insecure flag to zitadel.New)")
	key    = flag.String("key", "", "path to your key.json")
	port   = flag.String("port", "8089", "port to run the server on (default is 8089)")

	// tasks are used to store an in-memory list used in the protected endpoint
	tasks []string
)

/*
 This example demonstrates how to secure an HTTP API with ZITADEL using the provided authorization (AuthZ) middleware.

 It will serve the following 3 different endpoints:
 (These are meant to demonstrate the possibilities and do not follow REST best practices):

 - /api/healthz (can be called by anyone)
 - /api/tasks (requires authorization)
 - /api/add-task (requires authorization with granted `admin` role)
*/

func main() {
	flag.Parse()
	ctx := context.Background()

	// Initiate the authorization by providing a zitadel configuration and a verifier.
	// This example will use OAuth2 Introspection for this, therefore you will also need to provide the downloaded api key.json
	zit := zitadel.New(*domain, zitadel.WithInsecure("8080"))
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return
	}
	httpClient := httpClientNatsEspionage(nc)
	authZ, err := authorization.New(ctx, zit, DefaultAuthorization(*key, rs.WithClient(httpClient)))
	if err != nil {
		slog.Error("zitadel sdk could not initialize", "error", err)
		os.Exit(1)
	}
	mw := middleware.New(authZ)
	router := http.NewServeMux()

	// This endpoint is accessible by anyone and will always return "200 OK" to indicate the API is running
	router.Handle("/api/healthz", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err = jsonResponse(w, "OK", http.StatusOK)
			if err != nil {
				slog.Error("error writing response", "error", err)
			}
		}))

	// This endpoint is only accessible with a valid authorization (in this case a valid access_token / PAT).
	// It will list all stored tasks. In case the user is granted the `admin` role it will add a separate task telling him
	// to add a new task.
	router.Handle("/api/tasks", mw.RequireAuthorization()(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Using the [middleware.Context] function we can gather information about the authorized user.
			// This example will just print the users ID using the provided method, and it will also
			// print the username by directly access the field of the typed [*oauth.IntrospectionContext].
			authCtx := mw.Context(r.Context())
			slog.Info("user accessed task list", "id", authCtx.UserID(), "username", authCtx.Username)

			// Although this endpoint is accessible by any authorized user, you might want to take additional steps
			// if the user is granted a specific role. In this case an `admin` will be informed to add a new task:
			list := tasks
			slog.Info("user has roles", "roles", authCtx.Scope)
			if authCtx.IsGrantedRole("map-cli.write") {
				list = append(list, "create a new task on /api/add-task")
			}

			// return the existing task list
			err = jsonResponse(w, &taskList{Tasks: list}, http.StatusOK)
			if err != nil {
				slog.Error("error writing response", "error", err)
			}
		})))

	// This endpoint is only accessible with a valid authorization, which was granted the `admin` role (in any organization).
	// It will add the provided task to the list of existing ones.
	router.Handle("/api/add-task", mw.RequireAuthorization(authorization.WithRole(`admin`))(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// get the provided task and do not accept an empty value
			task := strings.TrimSpace(r.FormValue("task"))
			if task == "" {
				err = jsonResponse(w, "task must not be empty", http.StatusBadRequest)
				if err != nil {
					slog.Error("error writing invalid task response", "error", err)
					return
				}
				return
			}

			// since it was not empty, let's add it to the existing list
			tasks = append(tasks, task)

			// since we only want the authorized userID and don't need any specific data, we can simply use [authorization.UserID]
			slog.Info("admin added task", "id", authorization.UserID(r.Context()), "task", task)

			// inform the admin about the successful addition
			err = jsonResponse(w, fmt.Sprintf("task `%s` added", task), http.StatusOK)
			if err != nil {
				slog.Error("error writing task added response", "error", err)
				return
			}
		})))

	// start the server on the specified port (default http://localhost:8089)
	lis := fmt.Sprintf(":%s", *port)
	slog.Info("server listening, press ctrl+c to stop", "addr", "http://localhost"+lis)
	err = http.ListenAndServe(lis, router)
	if !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server terminated", "error", err)
		os.Exit(1)
	}
}

func httpClientNatsEspionage(natsConn *nats.Conn) *http.Client {
	return &http.Client{
		Transport: FartTransport{http.DefaultTransport, natsConn},
	}
}

type FartTransport struct {
	Base     http.RoundTripper
	natsConn *nats.Conn
}

func (t FartTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// "h8s.http.GET.localhost._well-known.openid-configuration"
	natsSubject := "h8s" + "." + req.URL.Scheme + "." + req.Method + "." + "localhost" + mapPath(req.URL.Path)
	fmt.Printf("Request URL: %s\n natsSubject: %s\n", req.URL.String(), natsSubject)
	// if strings.Contains(req.URL.Path, "openid") {
	// 	t.natsConn.Request("h8s.http.GET.localhost.", []byte(""), 5*time.Second)
	// }
	// if strings.Contains(req.URL.Path, "authorize") {
	// 	return t.Base.RoundTrip(req)
	// }
	body := []byte{}
	var err error
	if req.Body != nil {
		body, err = io.ReadAll(req.Body)
		fmt.Printf("Request body: %s\n", string(body))
	}
	if err != nil {
		fmt.Println("Error reading request body:", err)
	}
	msg := nats.Msg{
		Subject: natsSubject,
		Data:    body,
		Header:  fromHttpHeaders(req.Header),
	}
	fmt.Printf("httpRequest: %v\n", req)
	fmt.Printf("NATS headers: %v\n", msg.Header)
	fmt.Printf("nats data: %s\n", string(msg.Data))
	reply, err := t.natsConn.RequestMsg(&msg, time.Second*5)
	if err != nil {
		fmt.Println("Error requesting NATS subject:", err)
	} else {
		fmt.Printf("NATS reply: %s\n", string(reply.Data))
	}
	return toHttpResponse(reply, req)
}

func toHttpResponse(reply *nats.Msg, req *http.Request) (*http.Response, error) {
	if reply == nil {
		return nil, errors.New("nats reply is nil")
	}
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     fromNatsHeader(reply.Header),
		Body:       io.NopCloser(strings.NewReader(string(reply.Data))),
		Request:    req,
	}
	return resp, nil
}

func fromNatsHeader(header nats.Header) http.Header {
	h := http.Header{}
	for k, v := range header {
		if len(v) > 0 {
			h.Add(k, v[0])
		}
	}
	return h
}

func fromHttpHeaders(header http.Header) nats.Header {
	h := nats.Header{}
	for k, v := range header {
		if len(v) > 0 {
			h.Add(k, v[0])
		}
	}
	return h
}

func mapPath(s string) string {
	s = strings.ReplaceAll(s, ".", "_")
	s = strings.ReplaceAll(s, "/", ".")
	return s
}

// jsonResponse is a simple helper function to return a proper JSON response
func jsonResponse(w http.ResponseWriter, resp any, status int) error {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

type taskList struct {
	Tasks []string `json:"tasks,omitempty"`
}

func DefaultAuthorization(path string, opt rs.Option) authorization.VerifierInitializer[*oauth.IntrospectionContext] {
	c, err := client2.ConfigFromKeyFile(path)
	if err != nil {
		return func(ctx context.Context, _ *zitadel.Zitadel) (authorization.Verifier[*oauth.IntrospectionContext], error) {
			return nil, err
		}
	}
	return oauth.WithIntrospection[*oauth.IntrospectionContext](JWTProfileIntrospectionAuthentication(c, opt))
}

func JWTProfileIntrospectionAuthentication(file *client2.KeyFile, opt rs.Option) oauth.IntrospectionAuthentication {
	return func(ctx context.Context, issuer string) (rs.ResourceServer, error) {
		return rs.NewResourceServerJWTProfile(ctx, issuer, file.ClientID, file.KeyID, []byte(file.Key), opt)
	}
}
