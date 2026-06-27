package config

import "github.com/spf13/viper"

func Load(path string) (AppConfig, error) {
	cfg := Default()
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return cfg, err
	}
	return cfg, v.Unmarshal(&cfg)
}
