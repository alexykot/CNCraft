package control

import (
	"encoding/base64"
	"strings"

	"github.com/google/uuid"
)

type ServerConf struct {
	ServerID  string `yaml:"server-id"` // ID of the current server. Set to random 16 char string by default.
	Brand     string // Server brand. Is always set to `CNCraft` and cannot be changed.
	IsCracked bool   `yaml:"is-cracked"` // if True - skip player authentication, connection encryption and compression. Set to False by default.

	World WorldConf `yaml:"world"` // Configuration of the world to load. Only one world per server at a time supported.

	DBURL string      `yaml:"db-url"` // URL of the postgres server
	Net   NetworkConf `yaml:"network"`
	Log   logLevels   `yaml:"log-levels"` // log level settings for various subsystems.

}

type NetworkConf struct {
	Host        string `yaml:"host"` // resolvable hostname/IP to bind to. Set to `localhost` by default.
	Port        int    `yaml:"port"` // TCP port to serve on. Set to 25566 by default.
	ZipTreshold int32  // size of packet in bytes from which to start compressing the packets. Cannot be set externally.
}

type WorldConf struct {
	// ID of the world to load. Must be a 36-char UUID string. Identifies a world saved in persistence, server will
	// fail to start if the world with this ID is not found. World must be pre-created separately using migration tools.
	WorldID uuid.UUID `yaml:"world-id"`

	// World shard size in *chunks*. Will be used as XxZ size of the shard, i.e. will be 10x10 chunks per shard
	// if set to 10 (default). Min 1, max 64.
	ShardSize           int  `yaml:"shard-size"`
	EnableRespawnScreen bool // Enable respawn screen or tell client to respawn immediately.
}

// One of DEBUG, INFO, WARN, ERROR for each of the items. Set to `INFO` by default.
type logLevels struct {
	// Baseline defines minimum global log level. Any other log level may be more verbose, but cannot be less verbose.
	// Defaults to `ERROR`
	// DEBT baseline is in fact the opposite of the above - it's the *most* verbose level, not *least* verbose.
	//  `zap` only allows this mechanics, will need to either find a workaround for this issue or replace `zap`.
	Baseline string `yaml:"base"`

	Dispatcher string `yaml:"dispatcher"`
	Network    string `yaml:"network"`
	PubSub     string `yaml:"pubsub"`
	DB         string `yaml:"db"`

	Players string `yaml:"players"`
	Windows string `yaml:"windows"`
	World   string `yaml:"world"`
	Sharder string `yaml:"sharder"`
}

// currentConf is an internal singleton of server configuration. It is registered once during server bootstrap.
// It should not be changed in runtime, effects of the runtime change are undefined.
var currentConf ServerConf

// RegisterCurrentConfig registers a config as the singleton of configuration for current server process.
func RegisterCurrentConfig(serverConfig ServerConf) {
	currentConf = serverConfig
}

// GetCurrentConfig provides current singleton config.
func GetCurrentConfig() ServerConf {
	return currentConf
}

// GetDefaultConfig is a temporary function to provide a workable hardcoded development config. Eventually it will be
// replaced with a proper external configuration loader.
func GetDefaultConfig() ServerConf {
	return addDefaults(ServerConf{
		ServerID:  strings.ToUpper(base64.StdEncoding.EncodeToString([]byte(uuid.New().String())))[:16],
		Brand:     "CNCraft",
		IsCracked: true,

		World: WorldConf{
			WorldID:             uuid.New(),
			ShardSize:           3,
			EnableRespawnScreen: true,
		},

		DBURL: "postgresql://postgres:root@127.0.0.1:5432/cncraft?sslmode=disable",
		Net: NetworkConf{
			Host:        "0.0.0.0",
			Port:        25566,
			ZipTreshold: 1,
		},
		Log: logLevels{
			Baseline: "DEBUG",

			Dispatcher: "ERROR",
			Network:    "ERROR",
			PubSub:     "ERROR",
			DB:         "ERROR",

			Players: "ERROR",
			Windows: "ERROR",
			Sharder: "DEBUG",
			World:   "DEBUG",
		},
	})
}

func addDefaults(conf ServerConf) ServerConf {
	if conf.World.ShardSize < 1 || conf.World.ShardSize > 64 {
		conf.World.ShardSize = 10
	}

	if conf.Log.Baseline == "" {
		conf.Log.Baseline = "ERROR"
	}

	if conf.Log.Network == "" {
		conf.Log.Network = conf.Log.Baseline
	}
	if conf.Log.PubSub == "" {
		conf.Log.PubSub = conf.Log.Baseline
	}
	if conf.Log.Players == "" {
		conf.Log.Players = conf.Log.Baseline
	}
	if conf.Log.Windows == "" {
		conf.Log.Windows = conf.Log.Baseline
	}
	if conf.Log.World == "" {
		conf.Log.World = conf.Log.Baseline
	}
	if conf.Log.Sharder == "" {
		conf.Log.Sharder = conf.Log.Baseline
	}
	if conf.Log.DB == "" {
		conf.Log.DB = conf.Log.Baseline
	}

	if conf.Net.Host == "" {
		conf.Net.Host = "localhost"
	}

	if conf.Net.Port == 0 {
		conf.Net.Port = 25566
	}

	return conf
}
