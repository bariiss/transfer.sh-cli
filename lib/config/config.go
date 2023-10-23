package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Path: lib/config/config.go

type Config struct {
	BaseURL, User, Pass string
}

// LoadConfig loads the config file
func LoadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".config", "transfersh-cli", ".config")
	content, err := os.ReadFile(configPath)
	if err != nil {
		return CreateConfig(configPath)
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) < 3 {
		return nil, fmt.Errorf("invalid config format")
	}

	return &Config{
		BaseURL: lines[0],
		User:    lines[1],
		Pass:    lines[2],
	}, nil
}

// CreateConfig creates a new config file
func CreateConfig(path string) (*Config, error) {
	dir := filepath.Dir(path)                      // get directory
	if err := os.MkdirAll(dir, 0755); err != nil { // create directory
		return nil, err
	}

	var config Config // create config
	fmt.Print("Enter transfer base URL: ")
	if _, err := fmt.Scanln(&config.BaseURL); err != nil { // get base URL
		return nil, err
	}
	fmt.Print("Enter transfer user: ")
	if _, err := fmt.Scanln(&config.User); err != nil { // get user
		return nil, err
	}
	fmt.Print("Enter transfer pass: ")
	if _, err := fmt.Scanln(&config.Pass); err != nil { // get pass
		return nil, err
	}

	content := fmt.Sprintf("%s\n%s\n%s", config.BaseURL, config.User, config.Pass) // create content
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {              // write content
		return nil, err
	}

	return &config, nil
}
