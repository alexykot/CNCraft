package control

import (
	"encoding/base64"
	"strings"

	"github.com/google/uuid"
)

type ServerConf struct {
	DBURL               string `yaml:"db-url"` // URL of the postgres server
	Network             NetworkConf
	LogLevels           logLevels `yaml:"log-levels"` // log level settings for various subsystems.
	IsCracked           bool      `yaml:"is-cracked"` // if True - skip player authentication, connection encryption and compression. Set to False by default.
	ServerID            string    `yaml:"server-id"`  // ID of the current server. Set to random 16 char string by default.
	EnableRespawnScreen bool      // Enable respawn screen or tell client to respawn immediately.
	Brand               string    // Server brand. Is always set to `CNCraft`.
}

type NetworkConf struct {
	Host        string `yaml:"host"` // resolvable hostname/IP to bind to. Set to `localhost` by default.
	Port        int    `yaml:"port"` // TCP port to serve on. Set to 25566 by default.
	ZipTreshold int32  // size of packet in bytes from which to start compressing the packets. Cannot be set externally.
}

// One of DEBUG, INFO, WARN, ERROR for each of the items. Set to `INFO` by default.
type logLevels struct {
	// Baseline defines minimum global log level. Any other log level may be more verbose, but cannot be less verbose.
	// Defaults to `ERROR`
	// DEBT baseline is in fact the opposite of the above - it's the *most* verbose level, not *least* verbose.
	//  `zap` only allows this mechanics, will need to either workaround this issue or replace `zap`.
	Baseline string `yaml:"base"`

	Dispatcher string `yaml:"dispatcher"`
	Network    string `yaml:"network"`
	PubSub     string `yaml:"pubsub"`
	Players    string `yaml:"players"`
	Windows    string `yaml:"windows"`
	DB         string `yaml:"db"`
}

// currentConf is an internal singleton of server configuration. It is registered once during server bootstrap.
// It should not be changed in runtime, effects of the runtime change are undefined.
var currentConf ServerConf

func RegisterCurrentConfig(serverConfig ServerConf) {
	currentConf = serverConfig
}

func GetCurrentConfig() ServerConf {
	return currentConf
}

func GetDefaultConfig() ServerConf {
	return addDefaults(ServerConf{
		DBURL: "postgresql://postgres:root@127.0.0.1:5432/cncraft?sslmode=disable",
		Network: NetworkConf{
			Host:        "0.0.0.0",
			Port:        25565,
			ZipTreshold: 1,
		},
		LogLevels: logLevels{
			Baseline: "DEBUG",

			Dispatcher: "ERROR",
			Network:    "DEBUG",
			PubSub:     "ERROR",

			Players: "DEBUG",
			Windows: "DEBUG",
			DB:      "ERROR",
		},
		IsCracked:           true,
		EnableRespawnScreen: true,
		ServerID:            strings.ToUpper(base64.StdEncoding.EncodeToString([]byte(uuid.New().String())))[:16],
		Brand:               "CNCraft",
	})
}

func addDefaults(conf ServerConf) ServerConf {
	if conf.LogLevels.Baseline == "" {
		conf.LogLevels.Baseline = "ERROR"
	}

	if conf.LogLevels.Network == "" {
		conf.LogLevels.Network = conf.LogLevels.Baseline
	}
	if conf.LogLevels.PubSub == "" {
		conf.LogLevels.PubSub = conf.LogLevels.Baseline
	}
	if conf.LogLevels.Players == "" {
		conf.LogLevels.Players = conf.LogLevels.Baseline
	}
	if conf.LogLevels.Windows == "" {
		conf.LogLevels.Windows = conf.LogLevels.Baseline
	}
	if conf.LogLevels.DB == "" {
		conf.LogLevels.DB = conf.LogLevels.Baseline
	}

	if conf.Network.Host == "" {
		conf.Network.Host = "localhost"
	}

	if conf.Network.Port == 0 {
		conf.Network.Port = 25566
	}

	return conf
}
