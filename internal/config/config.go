package config

const (
	ENV_FILE            = ".env"
	DEFAULT_CONFIG_FILE = "configs/default.yaml"
)

type Config struct {
	DB       DBConfig
	Openai   OpenaiConfig
	Server   ServerConfig
	LogLevel string `mapstructure:"log_level"`
}
