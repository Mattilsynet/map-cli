package component

import (
	"log"
	"testing"
)

func TestGenerateComponent(t *testing.T) {
	log.Println("hey")
	config := Config{
		Path: "../project/go.mod.cue",
	}
	GenerateComponent(config)
}
