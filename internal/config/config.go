package config

type AppConfig struct {
	Endpoint             string `mapstructure:"endpoint" yaml:"endpoint"`
	ACID                 string `mapstructure:"acid" yaml:"acid"`
	Username             string `mapstructure:"username" yaml:"username"`
	CheckIntervalSeconds int    `mapstructure:"checkIntervalSeconds" yaml:"checkIntervalSeconds"`
	AutoConnect          bool   `mapstructure:"autoConnect" yaml:"autoConnect"`
	AutoReconnect        bool   `mapstructure:"autoReconnect" yaml:"autoReconnect"`
	LogLevel             string `mapstructure:"logLevel" yaml:"logLevel"`
}

func Default() AppConfig {
	return AppConfig{
		Endpoint:             "http://192.168.112.30",
		ACID:                 "0",
		CheckIntervalSeconds: 60,
		AutoConnect:          true,
		AutoReconnect:        true,
		LogLevel:             "info",
	}
}
