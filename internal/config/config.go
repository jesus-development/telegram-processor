package config

const (
	ENV_FILE            = ".env"
	DEFAULT_CONFIG_FILE = "config.yaml"
)

type Config struct {
	DB     DBConfig
	Openai OpenaiConfig
	Server ServerConfig
}
