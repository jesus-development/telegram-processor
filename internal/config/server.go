package config

type (
	ServerConfig struct {
		GRPC GRPCConfig
		HTTP HTTPConfig
	}

	GRPCConfig struct {
		Host string
		Port int64
	}

	HTTPConfig struct {
		Host string
		Port int64
	}
)
