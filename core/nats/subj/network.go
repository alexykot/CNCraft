package subj

import "github.com/google/uuid"

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

// MkPlayerSpatialUpdate creates a subject name string for announcing player position updates.
//  This is sent every time player position changes.
func MkPlayerSpatialUpdate() string { return "players.update.spatial" }

// MkNewPlayerJoined creates a subject name string for announcing new players joining server for the first time.
//  This is sent every time player position changes.
func MkNewPlayerJoined() string { return "players.joined.first_time" }
