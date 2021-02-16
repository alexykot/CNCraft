package control

import (
	"encoding/base64"
	"strings"

	"github.com/google/uuid"
)

// currentConf is an internal singleton of server configuration. It is registered once during server bootstrap.
var currentConf ServerConf

func RegisterCurrentConfig(serverConfig ServerConf) {
	currentConf = serverConfig
}

func GetCurrentConfig() ServerConf {
	return currentConf
}

func GetDefaultConfig() ServerConf {
	return ServerConf{
		DBURL: "postgresql://postgres:root@127.0.0.1:5432/cncraft?sslmode=disable",
		Network: NetworkConf{
			Host:        "0.0.0.0",
			Port:        25565,
			ZipTreshold: 1,
		},
		LogLevel:            "DEBUG",
		IsCracked:           true,
		EnableRespawnScreen: true,
		ServerID:            strings.ToUpper(base64.StdEncoding.EncodeToString([]byte(uuid.New().String())))[:16],
		Brand:               "CNCraft",
	}
}

type ServerConf struct {
	DBURL               string `yaml:"db-url"` // URL of the postgres server
	Network             NetworkConf
	LogLevel            string `yaml:"log-level"`  // one of DEBUG, INFO, WARN, ERROR. Set to `INFO` by default.
	IsCracked           bool   `yaml:"is-cracked"` // if True - skip player authentication, connection encryption and compression. Set to False by default.
	ServerID            string `yaml:"server-id"`  // ID of the current server. Set to random 16 char string by default.
	EnableRespawnScreen bool   // Enable respawn screen or tell client to respawn immediately.
	Brand               string // Server brand. Is always set to `CNCraft`.
}

type NetworkConf struct {
	Host        string `yaml:"host"` // resolvable hostname/IP to bind to. Set to `localhost` by default.
	Port        int    `yaml:"port"` // TCP port to serve on. Set to 25566 by default.
	ZipTreshold int32  // size of packet in bytes from which to start compressing the packets. Cannot be set externally.
}
