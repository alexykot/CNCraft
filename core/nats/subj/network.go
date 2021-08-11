package subj

import (
	"github.com/google/uuid"
)

type Subj string

func (s Subj) String() string {
	return string(s)
}

// *** Connection related subjects ***

// MkConnReceive creates a subject name string for given connection ID for receiving server bound packets.
func MkConnReceive(connID uuid.UUID) Subj { return Subj("conn.receive." + connID.String()) }

// MkConnTransmit creates a subject name string for given connection ID for transmitting client bound packets.
func MkConnTransmit(connID uuid.UUID) Subj { return Subj("conn.transmit." + connID.String()) }

// MkConnStateChange creates a subject name string for given connection ID for handling connection state changes.
func MkConnStateChange(connID uuid.UUID) Subj { return Subj("conn.state." + connID.String()) }

// MkNewConn creates a subject name string for announcing new connections appearing.
func MkNewConn() Subj { return "conn.new" }

// MkConnClose creates a subject name string for announcing connections to be closed.
func MkConnClose() Subj { return "conn.close" }

// MkConnBroadcast creates a subject name string for announcing broadcast packets that need to be
// sent to all (or eventually some) players.
func MkConnBroadcast() Subj { return "conn.broadcast" }
