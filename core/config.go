package core

func DefaultConfig() ServerConfig {
	return ServerConfig{
		Network: Network{
			Host: "0.0.0.0",
			Port: 25565,
		},
	}
}

type ServerConfig struct {
	Network  Network
	LogLevel string `yaml:"log-level"`
}

type Network struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
