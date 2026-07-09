package postpass

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"database_name"`
}

type PostpassConfig struct {
	Database             DatabaseConfig `yaml:"database"`
	ListenPort           int            `yaml:"listen_port"`
	QuickMediumThreshold int            `yaml:"quick_medium_threshold"`
	MediumSlowThreshold  int            `yaml:"medium_slow_threshold"`
}

func DefaultConfig() PostpassConfig {
	return PostpassConfig{
		Database: DatabaseConfig{
			Host:         "localhost",
			Port:         5432,
			User:         "readonly",
			Password:     "readonly",
			DatabaseName: "gis",
		},
		ListenPort:           8081,
		QuickMediumThreshold: 150,
		MediumSlowThreshold:  150000,
	}
}

func defaultConfigPath() (string, error) {
	xdg_config_home := os.Getenv("XDG_CONFIG_HOME")
	if xdg_config_home != "" {
		return filepath.Join(xdg_config_home, "postpass.yaml"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ".config", "postpass.yaml"), nil
}

func LoadConfig(path *string) (PostpassConfig, error) {
	cfg := DefaultConfig()

	var configPath string
	isUsingDefaultPath := path == nil
	if isUsingDefaultPath {
		var err error
		configPath, err = defaultConfigPath()
		if err != nil {
			return cfg, err
		}
	} else {
		configPath = *path
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		// Only continue if no path was specified and the default file does not exist
		if !(isUsingDefaultPath && os.IsNotExist(err)) {
			return cfg, fmt.Errorf("error reading config file %s: %w", configPath, err)
		}

		// Write the default config to the default config path
		if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
			return cfg, fmt.Errorf("error creating config directory: %w", err)
		}
		out, err := yaml.Marshal(cfg)
		if err != nil {
			return cfg, fmt.Errorf("error serializing default config: %w", err)
		}
		if err := os.WriteFile(configPath, out, 0644); err != nil {
			return cfg, fmt.Errorf("error writing default config to %s: %w", configPath, err)
		}
		return cfg, nil
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("error parsing config file %s: %w", configPath, err)
	}

	return cfg, nil
}
