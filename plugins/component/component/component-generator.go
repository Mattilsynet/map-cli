package component

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func GenerateComponent(config Config) {
	/*  TODO:
	        2. generate go files according to selected capabilities, name of component and path to put them
	        3. Generate wit files
		4. Generate wadm files
		5. Generate wasmcloud.toml
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
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	log.Println(currentDir)
	goModFile, err := os.ReadFile("./../project/go.mod.cue")
	if err != nil {
		fmt.Errorf("error reading go.mod.cue file: %v", err)
		return
	}
	fmt.Printf("go.mod.cue file content: %s\n", string(goModFile))
}
