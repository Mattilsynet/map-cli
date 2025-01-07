package main

import (
	"fmt"
	"os"

	"github.com/Mattilsynet/map-cli/plugins/license/types"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "license",
		Short:   "License plugin",
		Aliases: []string{"l"},
	}
	// TODO: Improvement to add such that we can have different licences produced by arguments or flags
	create := &cobra.Command{
		Use:   "create",
		Short: "Create MIT license file called LICENCE. Remember to run inside your projects root directory. Set ENV variable MC_LICENSE_AUTHOR='<your-org/user-name>'. \n The MC_LICENSE_AUTHOR will be used as author of LICENSE",
		Run: func(cmd *cobra.Command, args []string) {
			createLicenceFile(types.MIT_LICENSE)
		},
	}
	rootCmd.AddCommand(create)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createLicenceFile(license types.LICENCE) {
	file, err := os.Create("LICENSE")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(license.GetLicense())
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
	fmt.Println(license.GetType() + " LICENSE file generated.")
}
