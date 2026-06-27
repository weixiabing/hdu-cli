package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

func Load(path string) (AppConfig, error) {
	cfg := Default()
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) || os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	if !v.IsSet("endpoint") && v.IsSet("net.endpoint") {
		cfg.Endpoint = v.GetString("net.endpoint")
	}
	if !v.IsSet("acid") && v.IsSet("net.acid") {
		cfg.ACID = v.GetString("net.acid")
	}
	if !v.IsSet("username") && v.IsSet("net.auth.username") {
		cfg.Username = v.GetString("net.auth.username")
	}

	return cfg, nil
}
