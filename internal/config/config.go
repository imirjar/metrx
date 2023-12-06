package config

type Config struct {
	Port int
}

func New() *Config {
	return &Config{8080}
}
