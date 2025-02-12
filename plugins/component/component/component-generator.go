package component

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/Mattilsynet/map-cli/plugins/component/project"
)

func GenerateApp(config Config) error {
	mapOfContent, err := ReadAllTemplateFiles(config, project.Templs)
	if err != nil {
		return err
	}
	err = GenerateFiles(config.Path, mapOfContent)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GenerateFiles(projectRootPath string, mapOfContent map[string]string) error {
	for path, content := range mapOfContent {
		// INFO: Only make files if content from templating is not empty
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

func ReadAllTemplateFiles(config Config, tmpls map[string]string) (map[string]string, error) {
	mapOfContent := make(map[string]string)
	for key, tmpl := range tmpls {
		txtFile, err := ExecuteTmplWithData(config, tmpl)
		if err != nil {
			log.Println("error reading file: ", key, " with error: ", err)
			return nil, err
		}
		mapOfContent[key] = txtFile
	}
	return mapOfContent, nil
}

func ExecuteTmplWithData(data interface{}, tmplContent string) (string, error) {
	tmpl, err := template.New("module").Parse(tmplContent)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
