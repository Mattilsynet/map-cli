package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Nats struct {
		CredentialFilePath string `mapstructure:"credential_file_path"`
		Session            string `mapstructure:"session"`
	}
}

func DefaultConfig() *Config {
	return &Config{
		Nats: struct {
			CredentialFilePath string "mapstructure:\"credential_file_path\""
			Session            string "mapstructure:\"session\""
		}{".", ""},
	}
}

// TODO: Make a lot of this code runnable from config object, to read the config from file and store to file
func main() {
	// create login command
	rootCmd := &cobra.Command{
		Use:     "mc-auth",
		Short:   "Auth plugin",
		Aliases: []string{"a", "auth"},
	}
	nats := &cobra.Command{
		Use:   "nats",
		Short: "Authenticate through nats context",
	}
	rootCmd.AddCommand(nats)
	nats.AddCommand(
		&cobra.Command{
			Use:   "login",
			Short: "Choose nats context, i.e., `nats context select`",
			Run: func(cmd *cobra.Command, args []string) {
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
				homeDir, err := os.UserHomeDir()
				if err != nil {
					fmt.Printf("Failed to get home directory: %v\n", err)
					return
				}

				// Define the target directory and file
				configDir := filepath.Join(homeDir, ".config", "map-cli")
				configFile := filepath.Join(configDir, "config")

				// Create the directory
				err = os.MkdirAll(configDir, 0o755) // 0755 gives read/write/execute permissions for the user, and read/execute for others
				if err != nil {
					fmt.Printf("Failed to create directory %s: %v\n", configDir, err)
					return
				}
				viper.SetConfigName("config")  // Name of the config file (without extension)
				viper.SetConfigType("toml")    // File format
				viper.AddConfigPath(configDir) // Path to look for the config file

				// Check if the config file exists
				if _, err := os.Stat(configFile); os.IsNotExist(err) {
					// Write the default configuration to the file
					err = viper.SafeWriteConfigAs(configFile)
					config := DefaultConfig()
					viper.Set("Nats.CredentialFilePath", config.Nats.CredentialFilePath)
					viper.Set("Nats.Session", uuid.NewString())
					if err != nil {
						fmt.Printf("Failed to write default config: %v\n", err)
						return
					}
					fmt.Printf("Default configuration created at: %s\n", configFile)
				} else {
					// Read the existing config file
					err = viper.ReadInConfig()
					if err != nil {
						fmt.Printf("Failed to read config file: %v\n", err)
						return
					}
				}

				viper.Set("nats.CredentialFilePath", path)
				viper.Set("nats.Session", uuid.NewString())
				err = viper.WriteConfig()
				if err != nil {
					fmt.Printf("Failed to update config file: %v\n", err)
					return
				}
			},
		})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
