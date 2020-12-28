package control

import "github.com/google/uuid"

func DefaultConfig() ServerConf {
	return ServerConf{
		Network: NetworkConf{
			Host:        "0.0.0.0",
			Port:        25565,
			ZipTreshold: 1,
		},
		LogLevel:            "DEBUG",
		IsCracked:           false,
		EnableRespawnScreen: true,
		ServerID:            uuid.New().String(),
	}
}

type ServerConf struct {
	Network             NetworkConf
	LogLevel            string `yaml:"log-level"` // one of DEBUG, INFO, WARN, ERROR. Set to `INFO` by default.
	IsCracked           bool   // if True - skip player authentication, connection encryption and compression. Set to False by default.
	ServerID            string // ID of the current server. Set to random UUID by default.
	EnableRespawnScreen bool
}

type NetworkConf struct {
	Host        string `yaml:"host"` // resolvable hostname/IP to bind to. Set to `localhost` by default.
	Port        int    `yaml:"port"` // TCP port to serve on. Set to 25566 by default.
	ZipTreshold int32  // size of packet in bytes from which to start compressing the packets. Cannot be set externally.
}
