package component

import (
	"slices"
	"strings"
)

type Config struct {
	Path, ComponentName, Repository string
	Capabilities                    []string
	ComponentNatsConn,              // these booleans should be deducted from Capabilities list above
	ComponentNatsJetstream,
	ComponentNatsKeyValue bool

	WitPackage string // can be deducted from Repository
	ImportNatsCoreWit,
	ExportNatsCoreWit,
	ImportNatsJetstreamWit,
	ExportNatsJetstreamWit,
	ImportNatsKvWit bool
}

type Opt func(*Config)

func NewConfig(path, componentName, repository string, capabilities []string, opts ...Opt) Config {
	config := Config{
		Path:          path,
		ComponentName: componentName,
		Repository:    repository,
		Capabilities:  capabilities,
	}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}

func WithComponentCode() Opt {
	return func(config *Config) {
		config.ComponentNatsConn = slices.Contains(config.Capabilities, "nats-core")
		config.ComponentNatsJetstream = slices.Contains(config.Capabilities, "nats-jetstream")
		config.ComponentNatsKeyValue = slices.Contains(config.Capabilities, "nats-kv")
		// TODO: add wit deduction or add it as an option
	}
}

func WithWitPackage() Opt {
	return func(config *Config) {
		config.WitPackage = deductWitPackage(config.Repository) + ":" + config.ComponentName
	}
}

func deductWitPackage(repository string) string {
	paths := strings.Split(repository, "/")
	return paths[1]
}
