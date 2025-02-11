package component

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/Mattilsynet/map-cli/plugins/component/project"
)

func GenerateApp(config Config) error {
	/*  TODO:
	        2. generate go files according to selected capabilities, name of component and path to put them
	        3. Generate wit files
		4. Generate wadm files
	        6. add sdk files from custom capabilities to go.mod
		7. Generate working <component-name>-component.go with implemented requirements according to capabilities, e.g., handle etc
		8. Generate README.md
		9. Generate LICENSE
	GenerateGoMod()
	GenerateToolsWithSDKs()
	GenerateComponent()
	GenerateWITWorld()
	GenerateWasmcloudToml()
	GenerateLocalWadm()
	GenerateReadme()
	GenerateGitIgnore()
	GenerateGithubWorkflow()
	*/
	mapOfContent, err := ReadAllTemplateFiles(config, project.Templs)
	if err != nil {
		return err
	}
	err = GenerateFiles(config.Path, mapOfContent)
	if err != nil {
		return err
	}

	return nil
}

func GenerateFiles(projectRootPath string, mapOfContent map[string]string) error {
	for path, content := range mapOfContent {
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
			return nil, err
		}
		mapOfContent[key] = txtFile
	}
	return mapOfContent, nil
}

func ExecuteTmplWithData(data interface{}, tmplFile string) (string, error) {
	tmpl, err := template.New("module").Parse(tmplFile)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
