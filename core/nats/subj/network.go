package subj

import "github.com/google/uuid"

// *** Connection related subjects ***
// MkConnReceive creates a subject name string for given connection ID for receiving server bound packets.
func MkConnReceive(connID uuid.UUID) string { return "conn." + connID.String() + ".receive" }

// MkConnTransmit creates a subject name string for given connection ID for sending client bound packets.
func MkConnTransmit(connID uuid.UUID) string { return "conn." + connID.String() + ".send" }

// MkConnStateChange creates a subject name string for given connection ID for handling connection state changes.
func MkConnStateChange(connID uuid.UUID) string { return "conn." + connID.String() + ".state" }

// MkNewConn creates a subject name string for announcing new connections appearing.
func MkNewConn() string { return "conn.new" }

// *** User related subjects ***
// MkJoinedPlayers creates a subject name string for announcing new users joining server after successful login.
func MkJoinedPlayers() string { return "users.joined" }
