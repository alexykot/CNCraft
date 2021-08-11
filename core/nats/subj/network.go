package subj

import (
	"github.com/google/uuid"
)

type Subj string

// *** Connection related subjects ***

// MkConnReceive creates a subject name string for given connection ID for receiving server bound packets.
func MkConnReceive(connID uuid.UUID) Subj { return Subj("conn.receive." + connID.String()) }

// MkConnTransmit creates a subject name string for given connection ID for transmitting client bound packets.
func MkConnTransmit(connID uuid.UUID) Subj { return Subj("conn.transmit." + connID.String()) }

// MkConnStateChange creates a subject name string for given connection ID for handling connection state changes.
func MkConnStateChange(connID uuid.UUID) Subj { return Subj("conn.state." + connID.String()) }

// MkNewConn creates a subject name string for announcing new connections appearing.
func MkNewConn() Subj { return "conn.new" }

// MkConnClose creates a subject name string for announcing connection to be closed.
func MkConnClose() Subj { return "conn.close" }

// DEBT There will be a need for global transmission channel, to broadcast world state updates, chat messages etc.
//  This will likely be done via separate broadcasting channel, and it will need to be handled by the broadcaster,
//  a new component in the network subsystem. Broadcaster of every node will unicast messages to every player connected
//  to that node. Eventually the broadcaster will become more clever and select a subset of players that actually need
//  the message (e.g. only those subscribed to relevant chat channels, or close enough to the updated chunk).
func MkConnBroadcast() Subj { return "conn.broadcast" }
