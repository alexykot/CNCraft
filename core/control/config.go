package control

func DefaultConfig() ServerConf {
	return ServerConf{
		Network: NetworkConf{
			Host: "0.0.0.0",
			Port: 25565,
		},
	}
}

type ServerConf struct {
	Network  NetworkConf
	LogLevel string `yaml:"log-level"`
}

type NetworkConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
