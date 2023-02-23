package utils

type Config struct {
	LogLevel string `toml:"log_level"`
}

func NewConfig() *Config {
	return &Config{
		LogLevel: "debug",
	}
}
