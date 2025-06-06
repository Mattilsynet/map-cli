package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Mattilsynet/map-cli/internal/config"
	"github.com/Mattilsynet/mapis/gen/go/command/v1"
	me_gcp "github.com/Mattilsynet/mapis/gen/go/managedgcpenvironment/v1"
	metav1 "github.com/Mattilsynet/mapis/gen/go/meta/v1"
	"github.com/Mattilsynet/mapis/gen/go/query/v1"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

const (
	ApiSubject = "map.%s"
)

// TODO: implement somewhere else the nats connection and sendToMapQueryApi, as it's more general than just managed-environment

type ManagedEnvironmentHandler struct {
	nc *nats.Conn
}

func New(nc *nats.Conn) *ManagedEnvironmentHandler {
	return &ManagedEnvironmentHandler{nc: nc}
}

func (ma *ManagedEnvironmentHandler) HandleCobraCommand(cmd *cobra.Command, args []string) error {
	// Now, send the message to your API
	subject := "map.result." + config.CurrentConfig.Nats.Session + ".>"
	var subscription *nats.Subscription
	var subErr error
	go func() {
		subscription, subErr = ma.nc.Subscribe(subject, func(msg *nats.Msg) {
			var qry query.Query
			if err := proto.Unmarshal(msg.Data, &qry); err == nil {
				bdy := qry.Status.TypePayload
				for _, bytes := range bdy {
					me := me_gcp.ManagedGcpEnvironment{}
					proto.Unmarshal(bytes, &me)
					bytes, err := json.MarshalIndent(&me, "", " ")
					if err != nil {
						fmt.Println("Error marshalling json")
					}
					fmt.Println(string(bytes))
				}
			} else {
				fmt.Println(string(msg.Data))
			}
		})
		if subErr != nil {
			log.Fatalf("Error subscribing to NATS: %v\n", subErr)
		}
	}()

	if cmd.Use == "get" {
		name := ""
		if len(args) == 1 {
			name = args[0]
		} else if len(args) > 1 {
			return fmt.Errorf("too many arguments, expected 1, got %d", len(args))
		}
		listOfMeta := make([]*metav1.ObjectMeta, 0)
		for _, arg := range args {
			listOfMeta = append(listOfMeta, &metav1.ObjectMeta{Name: arg})
		}
		qryErr := ma.sendToMapQueryApi(&query.Query{
			Type: &metav1.TypeMeta{
				Kind:       "Query",
				ApiVersion: "v1",
			},
			Metadata: &metav1.ObjectMeta{Name: "ManagedGcpEnvironment", ResourceVersion: uuid.NewString()},
			Spec: &query.QuerySpec{
				Action:  "get",
				Type:    &metav1.TypeMeta{Kind: "ManagedGcpEnvironment", ApiVersion: "v1"},
				Session: config.CurrentConfig.Nats.Session,
				QueryFilter: &query.QueryFilter{
					Name: name,
				},
				TypeMetadata: listOfMeta,
			},
		})
		if qryErr != nil {
			return qryErr
		}
		// INFO: APPLY!
	} else {
		if len(args) < 1 {
			return errors.New("no file provided")
		}

		filePath := args[0]

		data, err := readFileContent(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file '%s': %w", filePath, err)
		}

		var message *me_gcp.ManagedGcpEnvironment
		format := getFileFormat(filePath)
		switch format {
		case "json":
			message, err = unmarshalJSONToProto(data)
		case "yaml":
			message, err = unmarshalYAMLToProto(data)
		default:
			return fmt.Errorf("unsupported file format for file '%s'", filePath)
		}
		if err != nil {
			meEmpty := &me_gcp.ManagedGcpEnvironment{Type: &metav1.TypeMeta{Kind: "ManagedGcpEnvironment", ApiVersion: "v1"}, Metadata: &metav1.ObjectMeta{Name: "map-dev", ResourceVersion: uuid.NewString()}, Spec: &me_gcp.ManagedGcpEnvironmentSpec{}}
			meBytes, jsonMarshalIdentErr := json.MarshalIndent(meEmpty, "", " ")
			if jsonMarshalIdentErr != nil {
				return jsonMarshalIdentErr
			}
			return fmt.Errorf("failed to unmarshal data from file '%s' with error: %w \n valid json would be: \n%s", filePath, err, string(meBytes))
		}
		bytes, protoMarshalErr := message.MarshalVT()
		if protoMarshalErr != nil {
			return protoMarshalErr
		}
		sessionId := config.CurrentConfig.Nats.Session
		if sessionId == "" {
			return errors.New("no session id provided, try map-cli auth login")
		}
		cmdErr := ma.sendToMapCommandApi(&command.Command{
			Type: &metav1.TypeMeta{
				Kind:       "Command",
				ApiVersion: "v1",
			},
			Metadata: &metav1.ObjectMeta{Name: "Command", ResourceVersion: uuid.NewString()},
			Spec: &command.CommandSpec{
				Operation:   cmd.Use,
				Type:        &metav1.TypeMeta{Kind: "ManagedGcpEnvironment", ApiVersion: "v1"},
				TypePayload: bytes,
				SessionId:   sessionId,
			},
		})
		if cmdErr != nil {
			return cmdErr
		}
	}
	time.Sleep(2 * time.Second)
	subscription.Unsubscribe()
	return nil
}

func (ma *ManagedEnvironmentHandler) sendToMapQueryApi(qry *query.Query) error {
	queryBytes, protoMarshalErr := qry.MarshalVT()
	if protoMarshalErr != nil {
		return protoMarshalErr
	}
	_, natsRequestErr := ma.nc.Request(fmt.Sprintf(ApiSubject, "get"), queryBytes, time.Second*10)
	if natsRequestErr != nil {
		return natsRequestErr
	}
	return nil
}

func (ma *ManagedEnvironmentHandler) sendToMapCommandApi(cmd *command.Command) error {
	commandBytes, protoMarshalErr := cmd.MarshalVT()
	if protoMarshalErr != nil {
		return protoMarshalErr
	}
	fmt.Println("command bytes: " + string(commandBytes))
	apiOperation := cmd.Spec.Operation
	a, natsRequestErr := ma.nc.Request(fmt.Sprintf(ApiSubject, apiOperation), commandBytes, time.Second*10)
	if natsRequestErr != nil {
		return natsRequestErr
	}
	fmt.Println("Response: ", string(a.Data))
	return nil
}

func readFileContent(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

func getFileFormat(filePath string) string {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".yaml", ".yml":
		return "yaml"
	case ".json":
		return "json"
	default:
		return "unknown"
	}
}

func unmarshalJSONToProto(data []byte) (*me_gcp.ManagedGcpEnvironment, error) {
	message := &me_gcp.ManagedGcpEnvironment{}
	unmarshalOptions := protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}
	err := unmarshalOptions.Unmarshal(data, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func unmarshalYAMLToProto(data []byte) (*me_gcp.ManagedGcpEnvironment, error) {
	var yamlMap map[string]any
	err := yaml.Unmarshal(data, &yamlMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	jsonData, err := json.Marshal(yamlMap)
	if err != nil {
		return nil, fmt.Errorf("failed to convert YAML to JSON: %w", err)
	}

	return unmarshalJSONToProto(jsonData)
}
