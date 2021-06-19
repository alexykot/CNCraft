package subj

import (
	"github.com/google/uuid"
)

// *** Connection related subjects ***

// MkConnReceive creates a subject name string for given connection ID for receiving server bound packets.
func MkConnReceive(connID uuid.UUID) string { return "conn." + connID.String() + ".receive" }

// MkConnTransmit creates a subject name string for given connection ID for transmitting client bound packets.
func MkConnTransmit(connID uuid.UUID) string { return "conn." + connID.String() + ".transmit" }

// MkConnStateChange creates a subject name string for given connection ID for handling connection state changes.
func MkConnStateChange(connID uuid.UUID) string { return "conn." + connID.String() + ".state" }

// MkNewConn creates a subject name string for announcing new connections appearing.
func MkNewConn() string { return "conn.new" }

// MkConnClose creates a subject name string for announcing connection to be closed.
func MkConnClose() string { return "conn.close" }

// *** Player related subjects ***

// MkPlayerLoading creates a subject name string for announcing new users joining server.
//  This is send after successful login and triggers client world loading and player spawn.
func MkPlayerLoading() string { return "players.loading" }

// MkPlayerJoined creates a subject name string for announcing new players successfully joined server.
//  This is sent after the player has successfully spawned in the world.
func MkPlayerJoined() string { return "players.joined" }

// MkPlayerLeft creates a subject name string for announcing new players successfully joined server.
//  This is sent after the player has successfully spawned in the world.
func MkPlayerLeft() string { return "players.left" }

// MkNewPlayerJoined creates a subject name string for announcing new players joining server for the first time.
//  This is sent every time player position changes.
func MkNewPlayerJoined() string { return "players.joined.first_time" }

// MkPlayerSpatialUpdate creates a subject name string for announcing player position updates.
//  This is sent every time player position changes.
func MkPlayerSpatialUpdate() string { return "players.update.spatial" }

// MkPlayerInventoryUpdate creates a subject name string for announcing player inventory updates.
//  This is sent every time player inventory changes (including hotbar).
func MkPlayerInventoryUpdate() string { return "players.update.inventory" }

// *** World related subjects ***

// MkShardEvent creates a subject name string for shard events, i.e. all events that are processed
// by the event loop of the given shard.
func MkShardEvent(shardID string) string {
	return "world.event." + shardID
}

// DEBT There will be a need for global transmission channel as well, to broadcast world state updates, chat messages etc.
//  This will likely be done via separate broadcasting channel, and it will need to be handled by the broadcaster,
//  a new component in the network subsystem. Broadcaster of every node will unicast messages to every player connected
//  to that node. Eventually the broadcaster will become more clever and select a subset of players that actually need
//  the message (e.g. only those subscribed to relevant chat channels, or close enough to the updated chunk).
// func MkConnBroadcast() string { return "conn.broadcast" }
