// mc-me/main.go
package main

import (
	"testing"

	"github.com/Mattilsynet/map-cli/internal/config"
)

func Test_main(t *testing.T) {
	_ = config.CurrentConfig
	tests := []struct {
		name string
	}{
		{name: "Test main function"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
