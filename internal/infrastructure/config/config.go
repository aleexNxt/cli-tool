package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/aleexNxt/cli-tool/internal/domain"
)

const (
	configFileName = ".devtool.json"
)

// ConfigRepository implementiert domain.ConfigRepository
type ConfigRepository struct {
	configPath string
}

// NewConfigRepository erstellt ein neues ConfigRepository
func NewConfigRepository() *ConfigRepository {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	return &ConfigRepository{
		configPath: filepath.Join(homeDir, configFileName),
	}
}

// NewConfigRepositoryWithPath erstellt ein ConfigRepository mit custom path
func NewConfigRepositoryWithPath(path string) *ConfigRepository {
	return &ConfigRepository{
		configPath: path,
	}
}

// Load lädt die Konfiguration
func (r *ConfigRepository) Load() (*domain.Config, error) {
	// Check if config file exists
	if _, err := os.Stat(r.configPath); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return domain.NewConfig(), nil
	}

	// Read config file
	data, err := os.ReadFile(r.configPath)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var config domain.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Ensure environment map is initialized
	if config.Environment == nil {
		config.Environment = make(map[string]string)
	}

	return &config, nil
}

// Save speichert die Konfiguration
func (r *ConfigRepository) Save(config *domain.Config) error {
	// Convert to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(r.configPath, data, 0644)
}

// Get holt einen einzelnen Konfigurationswert
func (r *ConfigRepository) Get(key string) (string, error) {
	config, err := r.Load()
	if err != nil {
		return "", err
	}

	switch key {
	case "build_command":
		return config.BuildCommand, nil
	case "test_command":
		return config.TestCommand, nil
	case "deploy_command":
		return config.DeployCommand, nil
	case "lint_command":
		return config.LintCommand, nil
	case "docker_registry":
		return config.DockerRegistry, nil
	default:
		if val, ok := config.Environment[key]; ok {
			return val, nil
		}
		return "", nil
	}
}

// Set setzt einen einzelnen Konfigurationswert
func (r *ConfigRepository) Set(key, value string) error {
	config, err := r.Load()
	if err != nil {
		return err
	}

	switch key {
	case "build_command":
		config.BuildCommand = value
	case "test_command":
		config.TestCommand = value
	case "deploy_command":
		config.DeployCommand = value
	case "lint_command":
		config.LintCommand = value
	case "docker_registry":
		config.DockerRegistry = value
	default:
		if config.Environment == nil {
			config.Environment = make(map[string]string)
		}
		config.Environment[key] = value
	}

	return r.Save(config)
}
