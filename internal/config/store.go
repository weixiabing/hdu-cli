package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
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

func Save(path string, cfg AppConfig, password string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}

	data := map[string]any{
		"endpoint":             cfg.Endpoint,
		"acid":                 cfg.ACID,
		"username":             cfg.Username,
		"checkIntervalSeconds": cfg.CheckIntervalSeconds,
		"autoConnect":          cfg.AutoConnect,
		"autoReconnect":        cfg.AutoReconnect,
		"logLevel":             cfg.LogLevel,
		"net": map[string]any{
			"endpoint": cfg.Endpoint,
			"acid":     cfg.ACID,
			"auth": map[string]any{
				"username": cfg.Username,
				"password": password,
			},
		},
	}

	content, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(path, content, 0o600)
}
