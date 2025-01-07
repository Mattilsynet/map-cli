package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Mattilsynet/map-cli/plugins/gitignore/gitignore"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "gitignore",
		Short:   "Gitignore plugin",
		Aliases: []string{"gi"},
	}
	// TODO: Improvement to add such that we can have different licences produced by arguments or flags
	create := &cobra.Command{
		Use:   "create",
		Short: "Creates a .gitignore file used for Map oriented repositories",
		Run: func(cmd *cobra.Command, args []string) {
			createGitIgnoreFile()
		},
	}
	rootCmd.AddCommand(create)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createGitIgnoreFile() {
	gitIgnoreFile, err := os.Create(".gitignore")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer gitIgnoreFile.Close()

	_, err = gitIgnoreFile.WriteString(gitignore.GitignoreTemplate)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
	fmt.Println(".gitignore file generated.")
}
