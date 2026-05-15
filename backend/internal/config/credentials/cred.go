package credentials

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

const _defaultCredentialsPath = "credentials.yaml"

type Credentials struct {
	PostgresDSN string `mapstructure:"PostgresDSN" validate:"required"`
}

type GrpcClient struct {
	Address        string        `mapstructure:"Address" validate:"required"`
	MaxMessageSize int           `mapstructure:"MaxMessageSize"`
	RetryAttempts  int           `mapstructure:"RetryAttempts"`
	MinRetryTime   time.Duration `mapstructure:"MinRetryTime"`
	MaxRetryTime   time.Duration `mapstructure:"MaxRetryTime"`
}

func New() (*Credentials, error) {
	vp := viper.New()

	vp.SetConfigFile(_defaultCredentialsPath)
	if err := vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading credentials: %w", err)
	}

	var creds Credentials
	if err := vp.Unmarshal(&creds); err != nil {
		return nil, fmt.Errorf("unmarshal credentials: %w", err)
	}

	if err := validator.New().Struct(&creds); err != nil {
		return nil, fmt.Errorf("validate credentials: %w", err)
	}

	return &creds, nil
}
