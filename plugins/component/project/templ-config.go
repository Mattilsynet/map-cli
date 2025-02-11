package project

import _ "embed"

//go:embed go.mod.templ
var goModtempl string

//go:embed component.go.templ
var componentGotempl string

//go:embed world.wit.templ
var worldWittempl string

//go:embed wadm.yaml.templ
var wadmYamltempl string

//go:embed wasmcloud.toml.templ
var wasmcloudTomltempl string

//go:embed README.md.templ
var readmeMdtempl string

//go:embed local.wadm.yaml.templ
var localWadmYamltempl string

//go:embed tools.go.templ
var toolsGotempl string

var Templs = map[string]string{
	"go.mod":          goModtempl,
	"component.go":    componentGotempl,
	"wadm.yaml":       wadmYamltempl,
	"local.wadm.yaml": localWadmYamltempl,
	"tools.go":        toolsGotempl,
	"wasmcloud.toml":  wasmcloudTomltempl,
	"README.md":       readmeMdtempl,
	"wit/world.wit":   worldWittempl,
}
