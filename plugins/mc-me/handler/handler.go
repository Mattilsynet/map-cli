package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Mattilsynet/map-cli/internal/config"
	"github.com/Mattilsynet/map-types/gen/go/command/v1"
	managedenvironment_v1 "github.com/Mattilsynet/map-types/gen/go/managedenvironment/v1"
	metav1 "github.com/Mattilsynet/map-types/gen/go/meta/v1"
	"github.com/Mattilsynet/map-types/gen/go/query/v1"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

// TODO: implement somewhere else the nats connection and sendToMapQueryApi, as it's more general than just managed-environment
type ManagedEnvironmentHandler struct {
	nc *nats.Conn
}

func New(nc *nats.Conn) *ManagedEnvironmentHandler {
	return &ManagedEnvironmentHandler{nc: nc}
}

func (ma *ManagedEnvironmentHandler) HandleCobraCommand(cmd cobra.Command, args ...string) error {
	// Now, send the message to your API
	if cmd.Use == "get" {
		qryErr := ma.sendToMapQueryApi(&query.Query{})
		if qryErr != nil {
			return qryErr
		}
	} else {
		if len(args) < 1 {
			return errors.New("no file provided")
		}
		filePath := args[0]

		data, err := readFileContent(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file '%s': %w", filePath, err)
		}

		var message *managedenvironment_v1.ManagedEnvironment
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
			return fmt.Errorf("failed to unmarshal data from file '%s': %w", filePath, err)
		}
		bytes, protoMarshalErr := proto.Marshal(message)
		if protoMarshalErr != nil {
			return protoMarshalErr
		}
		cmdErr := ma.sendToMapCommandApi(&command.Command{
			Type: &metav1.TypeMeta{
				Kind:       "Command",
				ApiVersion: "v1",
			},
			Metadata: &metav1.ObjectMeta{Name: "ManagedEnvironment", ResourceVersion: uuid.NewString()},
			Spec: &command.CommandSpec{
				Operation:   cmd.Use,
				Type:        &metav1.TypeMeta{Kind: "ManagedEnvironment", ApiVersion: "v1"},
				TypePayload: bytes,
				SessionId:   config.CurrentConfig.Nats.Session,
			},
		})
		if cmdErr != nil {
			return cmdErr
		}
	}

	return nil
}

func (ma *ManagedEnvironmentHandler) sendToMapQueryApi(qry *query.Query) error {
	panic("not implemented") // TODO: Implement
}

func (ma *ManagedEnvironmentHandler) sendToMapCommandApi(cmd *command.Command) error {
	panic("not implemented") // TODO: Implement
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

func unmarshalJSONToProto(data []byte) (*managedenvironment_v1.ManagedEnvironment, error) {
	message := &managedenvironment_v1.ManagedEnvironment{}
	unmarshalOptions := protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}
	err := unmarshalOptions.Unmarshal(data, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func unmarshalYAMLToProto(data []byte) (*managedenvironment_v1.ManagedEnvironment, error) {
	var yamlMap map[string]interface{}
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
