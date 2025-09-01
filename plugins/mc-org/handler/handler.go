package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	org "github.com/Mattilsynet/mapis/gen/go/organization"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

type OrgHandler struct {
	bearerToken string
}

func (o *OrgHandler) HandleCobraCommand(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("no file provided")
	}

	filePath := args[0]

	data, err := readFileContent(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file '%s': %w", filePath, err)
	}
	format := getFileFormat(filePath)
	var message *org.Organization
	switch format {
	case "json":
		message, err = unmarshalJSONToProto(data)
	case "yaml":
		message, err = unmarshalYAMLToProto(data)
	default:
		return fmt.Errorf("unsupported file format: %s", format)
	}
	if err != nil {
		return err
	}
	fmt.Println(message)
	return nil
}

func New(bearerToken string) *OrgHandler {
	return &OrgHandler{bearerToken: bearerToken}
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

func unmarshalJSONToProto(data []byte) (*org.Organization, error) {
	message := &org.Organization{}
	unmarshalOptions := protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}
	err := unmarshalOptions.Unmarshal(data, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func unmarshalYAMLToProto(data []byte) (*org.Organization, error) {
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
