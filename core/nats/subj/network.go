package subj

import "github.com/google/uuid"

// MkConnReceive creates a subject name string for given connection ID for receiving server bound packets
func MkConnReceive(connID uuid.UUID) string {
	return "conn." + connID.String() + ".receive"
}

// MkConnSend creates a subject name string for given connection ID for sending client bound packets
func MkConnSend(connID uuid.UUID) string {
	return "conn." + connID.String() + ".send"
}

// MkConnStateChange creates a subject name string for given connection ID for handling connection state changes
func MkConnStateChange(connID uuid.UUID) string {
	return "conn." + connID.String() + ".state"
}

// MkNewConn creates a subject name string for announcing new connections appearing
func MkNewConn() string {
	return "conn.new"
}

// MkNewUser creates a subject name string for announcing new connections appearing
func MkNewUser() string {
	return "user.new"
}
