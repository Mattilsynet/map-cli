package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CurrentConfig *Config

var CmdConfig = &cobra.Command{
	Use:   "config",
	Short: "Config",
}

var CmdConfigSet = &cobra.Command{}

type Nats struct {
	CredentialFilePath string `mapstructure:"credentialfilepath"`
	Session            string `mapstructure:"session"`
	nc                 *nats.Conn
}
type Config struct {
	Nats          Nats
	Authenticated bool
	AzureTenantID string
	AzureClientID string
}

func DefaultConfig() *Config {
	return &Config{
		Nats: Nats{
			CredentialFilePath: ".", Session: "",
		},
		Authenticated: false,
	}
}

func (n *Nats) GetConnection() (*nats.Conn, error) {
	if n.nc == nil {
		// TODO: do this more dynamically with using the credentials file path if exists, and determine nats url from config
		url, credsFilePath, natsContextErr := extractNatsContext(n.CredentialFilePath)
		var nc *nats.Conn
		if natsContextErr != nil {
			var err error
			log.Println("No nats context found, connecting to default nats server")
			nc, err = nats.Connect("nats://localhost:4222")
			if err != nil {
				return nil, err
			}
			n.nc = nc
			return n.nc, nil

		}
		var err error
		if credsFilePath != "" {
			log.Println("user creds and url found, connecting")
			nc, err = nats.Connect(url, nats.UserCredentials(credsFilePath))
		} else {
			log.Println("only connecting to url, no nats creds found")
			nc, err = nats.Connect(url)
		}
		if err != nil {
			return nil, err
		}
		n.nc = nc
	}
	return n.nc, nil
}

func extractNatsContext(credentialsFilePath string) (url string, natsCredentialsFilePath string, contextErr error) {
	content, err := os.Open(credentialsFilePath)
	if err != nil {
		return "", "", err
	}
	type Config struct {
		URL   string `json:"url"`
		Creds string `json:"creds"`
	}

	// Parse the JSON
	var config Config
	if err := json.NewDecoder(content).Decode(&config); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return "", "", err
	}

	// Extract the variables
	// read file from path into string
	// extract nats url from url in json
	// extract creds from creds in json
	return config.URL, config.Creds, nil
}

func init() {
	CurrentConfig = loadConfigIntoViper()
}

func IsAuthenticated() bool {
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
