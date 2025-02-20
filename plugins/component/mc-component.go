package main

import (
	"fmt"
	"os"

	"github.com/Mattilsynet/map-cli/plugins/component/component-generator"
	"github.com/Mattilsynet/map-cli/plugins/component/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "component",
		Short:   "Component plugin",
		Aliases: []string{"c"},
	}
	generate := &cobra.Command{
		Use:     "generate",
		Short:   "Generate a WasmCloud component",
		Aliases: []string{"gen", "g"},
		Run: func(cmd *cobra.Command, args []string) {
			m, err := tea.NewProgram(prompt.New()).Run()
			if err != nil {
				fmt.Println("error starting prompt program:", err)
				os.Exit(1)
			}
			if model := m.(*prompt.Model); model.Finished {
				generateModel, err := component.NewModel(model.ResultConfig())
				if err != nil {
					fmt.Println("error intiation, file-generation:", err)
				}
				teaModel, err := tea.NewProgram(generateModel).Run()
				if err != nil {
					fmt.Println("error starting code generation program:", err)
				}
				if model := teaModel.(component.Model); model.Done {
					fmt.Print("\n")
					fmt.Println("Goto the newly created component and follow README.md for further assistance")
					fmt.Print("\n")
					fmt.Println("cd " + model.RootPath)
				}

			}
		},
	}
	rootCmd.AddCommand(generate)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
