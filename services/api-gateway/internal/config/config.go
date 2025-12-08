package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort      string `mapstructure:"SERVER_PORT"`
	AuthServiceAddr string `mapstructure:"AUTH_SERVICE_ADDR"`
	UserServiceAddr string `mapstructure:"USER_SERVICE_ADDR"`
	JWTSecret       string `mapstructure:"JWT_SECRET"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.AddConfigPath("./configs")
	viper.SetConfigName("gateway_config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		// It's okay if config file doesn't exist, we can rely on env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}
