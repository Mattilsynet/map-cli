package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var CurrentConfig *Config

type Config struct {
	Nats struct {
		CredentialFilePath string `mapstructure:"credentialfilepath"`
		Session            string `mapstructure:"session"`
	}
	Authenticated bool
}

func DefaultConfig() *Config {
	return &Config{
		Nats: struct {
			CredentialFilePath string "mapstructure:\"credentialfilepath\""
			Session            string "mapstructure:\"session\""
		}{
			CredentialFilePath: ".", Session: "",
		},
		Authenticated: false,
	}
}

func init() {
	CurrentConfig = loadConfigIntoViper()
}

func isAuthenticated() bool {
	return CurrentConfig.Authenticated
}

func loadConfigIntoViper() *Config {
	// Load the configuration from the file
	// Define the target directory and file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Failed to get home directory: %v\n", err)
		os.Exit(1)
	}

	configDir := filepath.Join(homeDir, ".config", "map-cli")
	configFile := filepath.Join(configDir, "config")

	// Create the directory
	err = os.MkdirAll(configDir, 0o755) // 0755 gives read/write/execute permissions for the user, and read/execute for others
	if err != nil {
		fmt.Printf("Failed to create directory %s: %v\n", configDir, err)
		os.Exit(1)
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
			os.Exit(1)
		}
		fmt.Printf("Default configuration created at: %s\n", configFile)
	} else {
		// Read the existing config file
		err = viper.ReadInConfig()
		if err != nil {
			fmt.Printf("Failed to read config file: %v\n", err)
			os.Exit(1)
		}
	}

	configInFile := &Config{}
	err = viper.Unmarshal(configInFile)
	if err != nil {
		fmt.Printf("Failed to unmarshal config: %v\n", err)
		os.Exit(1)
	}
	return configInFile
}
