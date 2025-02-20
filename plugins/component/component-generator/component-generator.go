package component

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"slices"
	"text/template"

	project "github.com/Mattilsynet/map-cli/plugins/component/component-template"
)

type PathContent struct {
	Path, Content string
}

func GetTemplate(path string) (string, error) {
	for key, tmpl := range project.Templs {
		if key == path {
			return tmpl, nil
		}
	}
	return "", fmt.Errorf("template not found for path: %s", path)
}

func GetPathContentList(config *Config) ([]PathContent, error) {
	setBools(config)
	config.WitComponentName = toKebabCase(config.ComponentName)
	config.WitPackage = deductWitPackage(config.Repository) + ":" + config.WitComponentName
	pathContent, err := ReadAllTemplateFiles(config, project.Templs)
	if err != nil {
		return nil, err
	}

	return pathContent, nil
}

func setBools(config *Config) {
	// TODO: Really unfortunate logic to have to know what the above layer does, i.e., free text strings conveyed from tui, should be structured in a middle mapper
	config.ImportNatsCoreWit = slices.Contains(config.Capabilities, "nats-core:publish")
	config.ExportNatsCoreWit = slices.Contains(config.Capabilities, "nats-core:subscription")
	config.ExportNatsCoreRequestReplyWit = slices.Contains(config.Capabilities, "nats-core-wit")
	config.ImportNatsJetstreamWit = slices.Contains(config.Capabilities, "nats-jetstream:publish")
	config.ExportNatsJetstreamWit = slices.Contains(config.Capabilities, "nats-jetstream:consumer")
	config.ImportNatsKvWit = slices.Contains(config.Capabilities, "nats-kv:key-value")

	config.ComponentNatsCore = config.ImportNatsCoreWit || config.ExportNatsCoreWit
	config.ComponentNatsJetstream = config.ImportNatsJetstreamWit || config.ExportNatsJetstreamWit
	config.ComponentNatsKeyValue = config.ImportNatsKvWit
}

func GenerateAndInstall(projectRootPath, path, content string) error {
	if content != "" && content != "\n" {
		fullPath := projectRootPath + "/" + path
		if err := os.MkdirAll(getDirFromPath(fullPath), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
			return fmt.Errorf("failed to write file %s: %v", fullPath, err)
		}
	}
	return nil
}

func GenerateFiles(projectRootPath string, mapOfContent map[string]string) error {
	for path, content := range mapOfContent {
		if content != "" && content != "\n" {
			fullPath := projectRootPath + "/" + path
			if err := os.MkdirAll(getDirFromPath(fullPath), 0o755); err != nil {
				log.Println(err)
				return err
			}
			if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
				return fmt.Errorf("failed to write file %s: %v", fullPath, err)
			}
		}
	}
	return nil
}

func getDirFromPath(filePath string) string {
	lastSlash := -1
	for i, char := range filePath {
		if char == '/' || char == '\\' { // Handle both Unix and Windows paths
			lastSlash = i
		}
	}
	if lastSlash == -1 {
		return "." // No slashes, so return current directory
	}
	return filePath[:lastSlash]
}

func ReadTemplateFile(config *Config, tmpl string) (string, error) {
	content, err := ExecuteTmplWithData(config, tmpl)
	if err != nil {
		return "", err
	}
	return content, nil
}

func ReadAllTemplateFiles(config *Config, tmpls map[string]string) ([]PathContent, error) {
	pathContentList := make([]PathContent, 0)
	for key, tmpl := range tmpls {
		txtFile, err := ExecuteTmplWithData(config, tmpl)
		if err != nil {
			log.Println("error reading file: ", key, " with error: ", err)
			return nil, err
		}
		pathContentList = append(pathContentList, PathContent{Path: key, Content: txtFile})
	}
	return pathContentList, nil
}

func ExecuteTmplWithData(config *Config, tmplContent string) (string, error) {
	setBools(config)
	config.WitComponentName = toKebabCase(config.ComponentName)
	config.WitPackage = deductWitPackage(config.Repository) + ":" + config.WitComponentName

	tmpl, err := template.New("module").Parse(tmplContent)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		return "", err
	}
	return buf.String(), nil
}
