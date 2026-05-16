package configuration

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"time"

	"github.com/spf13/viper"
)

const _defaultConfigPath = "config.yaml"

type Config struct {
	DelayDataGen time.Duration `mapstructure:"DelayDataGen" validate:"required"`
}

func New() (*Config, error) {
	vp := viper.New()

	vp.SetConfigFile(_defaultConfigPath)
	if err := vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var conf Config
	if err := vp.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := validator.New().Struct(&conf); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &conf, nil
}
