package main

import (
	"fmt"
	"os"

	"github.com/Mattilsynet/map-cli/plugins/component/component-generator"
	"github.com/Mattilsynet/map-cli/plugins/component/tea-model"
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
			m, err := tea.NewProgram(teaModel.New()).Run()
			if err != nil {
				fmt.Println("error starting program:", err)
				os.Exit(1)
			}
			if model := m.(*teaModel.Model); model.Finished {
				component.GenerateApp(model.ResultConfig())
				/*  TODO:
				Feature: Add fancy loading bar
				Feature: Add fancy display of files generated in which folder
				Feature: cd to the newly creatid path
				*/
			}
		},
	}
	rootCmd.AddCommand(generate)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
