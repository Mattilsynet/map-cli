package component

import (
	"slices"
	"strings"
)

type Config struct {
	Path, ComponentName, Repository string
	Capabilities                    []string
	ComponentNatsCore,              // these booleans should be deducted from Capabilities list above
	ComponentNatsJetstream,
	ComponentNatsKeyValue bool

	WitPackage       string // can be deducted from Repository
	WitComponentName string // kebab-case wit requirement
	ImportNatsCoreWit,
	ExportNatsCoreWit,
	ExportNatsCoreRequestReplyWit,
	ImportNatsJetstreamWit,
	ExportNatsJetstreamWit,
	ImportNatsKvWit,
	ExportNatsKvWit bool
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
		config.ComponentNatsCore = slices.Contains(config.Capabilities, "nats-core")
		config.ComponentNatsJetstream = slices.Contains(config.Capabilities, "nats-jetstream")
		config.ComponentNatsKeyValue = slices.Contains(config.Capabilities, "nats-kv")
	}
}

func WithWit() Opt {
	return func(config *Config) {
		config.WitComponentName = toKebabCase(config.ComponentName)
		config.WitPackage = deductWitPackage(config.Repository) + ":" + config.WitComponentName
	}
}

func deductWitPackage(repository string) string {
	paths := strings.Split(repository, "/")
	return strings.ToLower(paths[1])
}

func toKebabCase(input string) string {
	var result strings.Builder
	for i, char := range input {
		if char >= 'A' && char <= 'Z' {
			// If it's not the first character, prepend a dash
			if i > 0 {
				result.WriteRune('-')
			}
			// Convert the uppercase letter to lowercase
			result.WriteRune(rune(strings.ToLower(string(char))[0]))
		} else {
			// Append non-uppercase characters as-is
			result.WriteRune(char)
		}
	}
	return result.String()
}
