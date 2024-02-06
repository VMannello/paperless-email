package pmail

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the overall configuration structure.
type Config struct {
	Accounts   []EmailAccount     `yaml:"accounts"`
	TagMapping map[string]Message `yaml:"tags"`
}

// LoadConfig reads the configuration from the specified YAML file.
func LoadConfig(filePath string) (Config, error) {
	var cfg Config

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(fileContent, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
