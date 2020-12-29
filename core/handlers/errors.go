package handlers

import (
	"fmt"
	"strings"
)

type PacketError struct {
	wrapped error
	error   errType
}

func (e PacketError) Error() string { return string(e.error) }
func (e PacketError) Unwrap() error { return e.wrapped }
func (e PacketError) Is(target error) bool {
	if targetErr, ok := target.(errType); ok {
		return strings.HasPrefix(e.Error(), targetErr.Error())
	}
	return false
}

type errType string

func (e errType) Error() string { return string(e) }

// List of sentinel errors
const InvalidLoginErr errType = "user login data invalid"

func newPacketError(topErr errType, wrappedErr error) PacketError {
	wrappedMessage := fmt.Sprintf("%s: %s", topErr, wrappedErr.Error())
	return PacketError{error: errType(wrappedMessage), wrapped: wrappedErr}
}
