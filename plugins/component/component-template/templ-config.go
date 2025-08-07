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

//go:embed pkgnats.go.templ
var pkgNatsGoTempl string

//go:embed pkgnatskv.go.templ
var pkgNatsKvTempl string

//go:embed pkgnatsjs.go.templ
var pkgNatsJsTempl string

//go:embed deps.toml.templ
var depsTomlTempl string

//go:embed github.workflows.yml.templ
var githubWorkflowsYmlTempl string

//go:embed pkgcronjob.go.templ
var pkgCronJobTempl string

//go:embed pkgcloudrunjob-admin-provider.map-me-gcp-cloudrunjob.go.templ
var pkgCloudrunjobAdminProviderMapMeGcpCloudrunJobTempl string

//go:embed pkgmanaged-environment.managed-environment.go.templ
var pkgManagedEnvironmentManagedEnvironmentTempl string

//go:embed pkgmanifest.manifest.go.templ
var pkgManifestManifestGoTempl string

const (
	ComponentGoPath   = "component.go"
	LocalWadmYamlPath = "local.wadm.yaml"
	GoModPath         = "go.mod"
	WadmYamlPath      = "wadm.yaml"
	ReadmeMdPath      = "README.md"
)

var Templs = map[string]string{
	GoModPath:                          goModtempl,
	ComponentGoPath:                    componentGotempl,
	WadmYamlPath:                       wadmYamltempl,
	LocalWadmYamlPath:                  localWadmYamltempl,
	ReadmeMdPath:                       readmeMdtempl,
	".github/workflows/build_push.yml": githubWorkflowsYmlTempl,
	"tools.go":                         toolsGotempl,
	"wasmcloud.toml":                   wasmcloudTomltempl,
	"wit/deps.toml":                    depsTomlTempl,
	"wit/world.wit":                    worldWittempl,
	// nats
	"pkg/nats/nats.go": pkgNatsGoTempl,
	"pkg/nats/kv.go":   pkgNatsKvTempl,
	"pkg/nats/js.go":   pkgNatsJsTempl,
	// cronjob
	"pkg/cronjob/cronjob.go": pkgCronJobTempl,
	// cloudrunjob-admin
	"pkg/cloudrunjob-admin-provider/map-me-gcp-cloudrunjob.go": pkgCloudrunjobAdminProviderMapMeGcpCloudrunJobTempl,
	"pkg/managed-environment/managed-environment.go":           pkgManagedEnvironmentManagedEnvironmentTempl,
	"pkg/manifest/manifest.go":                                 pkgManifestManifestGoTempl,
}
