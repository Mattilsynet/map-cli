package main

import (
	"testing"

	"github.com/Mattilsynet/map-cli/internal/config"
)

func Test_execLogin(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test execLogin",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execLogin(config.CurrentConfig.Nats.CredentialFilePath)
		})
	}
}
