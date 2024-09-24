package bootstrap

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// Config struct defines the configuration items for the application
type Config struct {
	// Server configuration
	ServerPort int `mapstructure:"SERVER_PORT" validate:"required"`

	// Database configuration
	DBDSN        string `mapstructure:"DB_DSN" validate:"required"`
	DBReplicaDSN string `mapstructure:"DB_REPLICA_DSN"`

	// SQS configuration
	SQSQueueURL            string `mapstructure:"SQS_QUEUE_URL" validate:"required,url"`
	SQSRegion              string `mapstructure:"SQS_REGION" validate:"required"`
	SQSMaxNumberOfMessages int32  `mapstructure:"SQS_MAX_NUMBER_OF_MESSAGES" validate:"required,min=1,max=10"`
	SQSWaitTimeSeconds     int32  `mapstructure:"SQS_WAIT_TIME_SECONDS" validate:"required,min=0,max=20"`
	SQSVisibilityTimeout   int32  `mapstructure:"SQS_VISIBILITY_TIMEOUT" validate:"required,min=0"`

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
	v.AutomaticEnv() // read env

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	if err := validator.New().Struct(config); err != nil {
		return nil, err
	}

	return &config, nil
}
