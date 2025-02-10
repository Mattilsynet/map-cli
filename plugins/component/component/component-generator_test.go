package component

import (
	"testing"
)

func TestGenerateComponent(t *testing.T) {
	config := Config{
		Path:          "/home/solve/git/temp",
		ComponentName: "test-component",
		Repository:    "github.com/Mattilsynet/test-component",
		Capabilities:  []string{"nats-core", "nats-jetstream"},
	}
	GenerateApp(config)
}
