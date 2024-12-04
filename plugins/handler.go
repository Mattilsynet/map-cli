package plugins

import (
	"github.com/spf13/cobra"
)

type Handler interface {
	HandleCobraCommand(cmd cobra.Command, args ...string) error
}
