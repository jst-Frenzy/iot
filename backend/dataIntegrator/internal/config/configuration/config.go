package configuration

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/spf13/viper"
)

const _defaultConfigPath = "config.yaml"

type Config struct {
	GRPCServer *GRPCServer `mapstructure:"GRPCServer" validate:"required"`
	WsServer   *WsServer   `mapstructure:"WsServer" validate:"required"`
	HTTPServer *HTTPServer `mapstructure:"HTTPServer" validate:"required"`
}

type WsServer struct {
	Address string `mapstructure:"Address" validate:"required"`
}

type GRPCServer struct {
	Address           string `mapstructure:"Address"           validate:"required"`
	MaxReceiveMsgSize int    `mapstructure:"MaxReceiveMsgSize" validate:"gte=0"`
	MaxSendMsgSize    int    `mapstructure:"MaxSendMsgSize"    validate:"gte=0"`
}

type HTTPServer struct {
	Address string `mapstructure:"Address" validate:"required"`
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
