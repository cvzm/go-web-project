package bootstrap

import (
	"github.com/spf13/viper"
)

// Config struct defines the configuration items for the application
type Config struct {
	// Server configuration
	ServerPort int `mapstructure:"SERVER_PORT"`

	// Database configuration
	DBDSN        string `mapstructure:"DB_DSN"`
	DBReplicaDSN string `mapstructure:"DB_REPLICA_DSN"`

	// External service configuration
	SlackWebhook string `mapstructure:"SLACK_WEBHOOK"`
}

// Constants related to configuration
const (
	EnvFile = ".env" // Environment file name
)

// NewConfig loads the configuration from env file and environment variables
func NewConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(EnvFile)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
