package main

import (
	"fmt"
	"os"

	"github.com/Mattilsynet/map-cli/plugins/component/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "provider",
		Short:   "Provider plugin",
		Aliases: []string{"p"},
	}
	generate := &cobra.Command{
		Use:     "generate",
		Short:   "Generate a WasmCloud provider",
		Aliases: []string{"gen", "g"},
		Run: func(cmd *cobra.Command, args []string) {
			promptModel, err := prompt.New()
			if err != nil {
				return
			}
			m, err := tea.NewProgram(promptModel).Run()
			if err != nil {
				fmt.Println("error starting prompt program:", err)
				os.Exit(1)
			}
			if model := m.(*prompt.Model); model.Finished {
			}
		},
	}
	rootCmd.AddCommand(generate)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
