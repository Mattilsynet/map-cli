module github.com/Mattilsynet/map-cli

go 1.23.0

toolchain go1.23.5

require (
	github.com/Mattilsynet/mapis v0.0.1
	github.com/charmbracelet/bubbles v0.20.0
	github.com/charmbracelet/bubbletea v1.1.0
	github.com/charmbracelet/lipgloss v1.0.0
	github.com/google/uuid v1.6.0
	github.com/nats-io/nats.go v1.38.0
	github.com/spf13/cobra v1.8.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.19.0
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9
	golang.org/x/oauth2 v0.18.0
	google.golang.org/api v0.171.0
	google.golang.org/protobuf v1.36.2
	gopkg.in/yaml.v3 v3.0.1
)

replace golang.org/x/crypto => golang.org/x/crypto v0.31.0
