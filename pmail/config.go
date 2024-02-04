package pmail

import (
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

// SMTPConfig represents SMTP configuration for an email account.
type SMTPConfig struct {
	Server             string `yaml:"server"`
	Port               int    `yaml:"port"`
	Username           string `yaml:"username"`
	Password           string `yaml:"password"`
	EnableTLS          bool   `yaml:"tls"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify"`
}

// EmailAccount represents an email account configuration.
type EmailAccount struct {
	Name  string     `yaml:"name"`
	Email string     `yaml:"email"`
	SMTP  SMTPConfig `yaml:"smtp"`
}

// TagMapping represents the overall tag mapping configuration.
type TagMapping struct {
	Config map[string]TagMappingConfig `yaml:"mapping"`
}

// TagMappingConfig represents the configuration for tag mappings.
type TagMappingConfig struct {
	Account    string   `yaml:"account"`
	Recipients []string `yaml:"recipients"`
}

// Message represents the configuration for the email message.
type Message struct {
	Subject           string `yaml:"subject"`
	Body              string `yaml:"body"`
	IncludeAttachment bool   `yaml:"include_attachment"`
}

// Config represents the overall configuration structure.
type Config struct {
	Accounts   []EmailAccount `yaml:"accounts"`
	TagMapping TagMapping     `yaml:"tags"`
	Message    Message        `yaml:"message"`
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

	// Replace environment variable placeholders in the message body
	cfg.Message.Body = replaceEnvironmentVariables(cfg.Message.Body)
	cfg.Message.Subject = replaceEnvironmentVariables(cfg.Message.Subject)

	return cfg, nil
}

var curlyRe = regexp.MustCompile(`\{\{([^\}]+)\}\}`)

func replaceEnvironmentVariables(input string) string {
	return curlyRe.ReplaceAllStringFunc(input, func(match string) string {
		return os.Getenv(strings.TrimSpace(strings.Trim(match, "{}")))
	})
}
