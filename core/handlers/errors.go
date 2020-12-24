package handlers

type PacketError struct {
	wrapped error
	error   errType
}

func (e PacketError) Error() string { return string(e.error) }

func (e PacketError) Unwrap() error { return e.wrapped }
func (e PacketError) Is(target error) bool {
	if targetErr, ok := target.(errType); ok {
		return targetErr.Error() == e.Error()
	}
	return false
}

type errType string

func (e errType) Error() string { return string(e) }

// List of sentinel errors
const InvalidLoginErr errType = "user login data invalid"

func newPacketError(errType errType, wrapped error) PacketError {
	return PacketError{error: errType, wrapped: wrapped}
}
