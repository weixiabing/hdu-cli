package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func Load(path string) (AppConfig, error) {
	cfg := Default()
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return cfg, err
	}
	return cfg, v.Unmarshal(&cfg)
}

func Save(path string, cfg AppConfig) error {
	v := viper.New()
	v.SetConfigFile(path)
	v.Set("endpoint", cfg.Endpoint)
	v.Set("acid", cfg.ACID)
	v.Set("username", cfg.Username)
	v.Set("checkIntervalSeconds", cfg.CheckIntervalSeconds)
	v.Set("autoConnect", cfg.AutoConnect)
	v.Set("autoReconnect", cfg.AutoReconnect)
	v.Set("logLevel", cfg.LogLevel)

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return v.WriteConfigAs(path)
}
