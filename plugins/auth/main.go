package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/Mattilsynet/map-cli/internal/config"
	_ "github.com/Mattilsynet/map-cli/internal/logger"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:     "mc-auth",
	Short:   "Auth plugin",
	Aliases: []string{"a", "auth"},
}

var natsCmd = &cobra.Command{
	Use:   "nats",
	Short: "Authenticate through nats context",
}

var natsCmdLogin = &cobra.Command{
	Use:   "login",
	Short: "Choose nats context, i.e., `nats context select`",
	Run: func(cmd *cobra.Command, args []string) {
		natsContextSelect()
	},
}

func init() {
	rootCmd.AddCommand(natsCmd)
	natsCmd.AddCommand(natsCmdLogin)
	rootCmd.AddCommand(azureCmd)
	azureCmd.AddCommand(azureCmdLogin)
}

func main() {
	pflag.Parse()
	rootCmd.Flags().AddFlagSet(pflag.CommandLine)

	slog.Debug("auth plugin executing")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func natsContextSelect() {
	// TODO: execplugin, nats context select and execContextSelect isnt needed as nats context select returns the info we need, so we should do both in same operation for efficiency
	err := execPlugin("nats", "context", "select")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute plugin: %v\n", err)
		os.Exit(1)
	}
	path, err := execContextSelect()
	if err != nil {
		fmt.Printf("Failed to get context path: %v\n", err)
		os.Exit(1)
	}
	execLogin(path)
}

func execLogin(path string) error {
	// Read the existing config file
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)
		return err
	}

	configInFile := config.CurrentConfig
	err = viper.Unmarshal(configInFile)
	if err != nil {
		fmt.Printf("Failed to unmarshal config: %v\n", err)
		return err
	}
	if strings.Compare(configInFile.Nats.CredentialFilePath, path) != 0 {
		viper.Set("nats.CredentialFilePath", path)
		viper.Set("nats.Session", uuid.NewString())
		err = viper.WriteConfig()
		if err != nil {
			fmt.Printf("Failed to update config file: %v\n", err)
			return err
		}
	}
	return nil
}

func execContextSelect() (string, error) {
	path, err := exec.LookPath("nats")
	if err != nil {
		return "", err
	}
	cmd := exec.Command(path, "context", "info")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	output := stdout.String()
	path = extractField(output, `Path:\s+(.+)`)
	return path, nil
}

func extractField(output, regex string) string {
	re := regexp.MustCompile(regex)
	match := re.FindStringSubmatch(output)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

func execPlugin(pluginName string, args ...string) error {
	path, err := exec.LookPath(pluginName)
	if err != nil {
		return fmt.Errorf("plugin '%s' not found in PATH", pluginName)
	}

	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
