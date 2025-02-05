package prompt

import (
	"fmt"

	"github.com/Mattilsynet/map-cli/plugins/component/pkg/component"
)

func PromptUser() component.ComponentConfig {
	componentConfig := component.ComponentConfig{}
	fmt.Print("Enter the component name: ")
	fmt.Scanln(&componentConfig.ComponentName)
	fmt.Print("Enter the full-path to install component: ")
	fmt.Scanln(&componentConfig.RootPath)

	fmt.Print("Include nats-core? (y/n): ")

	fmt.Scanln(&componentConfig.NatsCore)
	// component-name
	// root-path
	// nats-core
	// nats-jetstream
	// nats-kv
	// license
	panic("unimplemented")
}
